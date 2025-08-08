package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	MinScore     int64  `koanf:"min_score"`
	DebounceTime int64  `koanf:"debounce_time"`
	DefaultScope string `koanf:"default_scope"`
	FormatterCmd string `koanf:"formatter_cmd"`

	Scopes map[string]Scope `koanf:"scopes"`

	// Origins of a set configuration value, used for tracking
	// when/where the most recent value of a field was set
	// for debugging configurations.
	fieldOrigins map[string]string
}

type Scope struct {
	Name            string `koanf:"-"`
	Description     string `koanf:"description"`
	OptionsListFile string `koanf:"options-list-file"`
	OptionsListCmd  string `koanf:"options-list-cmd"`
	EvaluatorCmd    string `koanf:"evaluator"`
}

func NewConfig() *Config {
	return &Config{
		MinScore:     1,
		DebounceTime: 25,
		FormatterCmd: "nixfmt",

		Scopes: make(map[string]Scope),
	}
}

func ParseConfig(location ...string) (*Config, error) {
	k := koanf.New(".")

	fieldOrigins := make(map[string]string)

	for _, loc := range location {
		if _, err := os.Stat(loc); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return nil, err
		}

		fileK := koanf.New(".")

		err := fileK.Load(file.Provider(loc), toml.Parser())
		if err != nil {
			return nil, err
		}

		for _, key := range fileK.Keys() {
			fieldOrigins[key] = loc
		}

		// Also load incomplete scope keys into the field origins, since
		// scopes without proper definitions can technically exist.
		if scopesMap, ok := fileK.Get("scopes").(map[string]interface{}); ok {
			for scopeName := range scopesMap {
				scopeKey := fmt.Sprintf("scopes.%s", scopeName)
				fieldOrigins[scopeKey] = loc
			}
		}

		if err := k.Merge(fileK); err != nil {
			return nil, err
		}

	}

	cfg := NewConfig()
	cfg.fieldOrigins = fieldOrigins

	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	for name, scope := range cfg.Scopes {
		scope.Name = name
		cfg.Scopes[name] = scope
	}

	return cfg, nil
}

var optionTemplateRegex = regexp.MustCompile(`{{\s*-?\s*\.Option\b[^}]*}}`)

type ValidationError struct {
	Msg    string
	Origin string
}

func (e ValidationError) Error() string {
	msg := e.Msg
	if e.Origin != "" {
		msg += "\n\n" + color.YellowString("hint: this setting was last defined in %v", e.Origin)
	}
	return msg
}

func (c *Config) Validate() error {
	for s, v := range c.Scopes {
		if v.OptionsListCmd == "" && v.OptionsListFile == "" {
			return ValidationError{
				Msg:    fmt.Sprintf("no option list source defined for scope '%v'", s),
				Origin: c.FieldOrigin(fmt.Sprintf("scopes.%v", s)),
			}
		}
	}

	if c.DefaultScope != "" {
		foundScope := false
		for n := range c.Scopes {
			if n == c.DefaultScope {
				foundScope = true
				break
			}
		}

		if !foundScope {
			return ValidationError{
				Msg:    fmt.Sprintf("default scope '%v' not found", c.DefaultScope),
				Origin: c.FieldOrigin("default_scope"),
			}
		}
	}

	for s, v := range c.Scopes {
		if v.EvaluatorCmd == "" {
			continue
		}

		matches := optionTemplateRegex.FindAllString(v.EvaluatorCmd, -1)
		if len(matches) != 1 {
			origin := c.FieldOrigin(fmt.Sprintf("scopes.%v.evaluator", s))
			if len(matches) == 0 {
				return ValidationError{
					Msg:    fmt.Sprintf("evaluator for scope '%v' does not contain the placeholder {{ .Option }}", s),
					Origin: origin,
				}
			} else {
				return ValidationError{
					Msg:    fmt.Sprintf("multiple instances of {{ .Option }} placeholder in evaluator for scope '%v'", s),
					Origin: origin,
				}
			}
		}
	}

	return nil
}

func (c *Config) FieldOrigin(key string) string {
	if c.fieldOrigins == nil {
		return ""
	}

	return c.fieldOrigins[key]
}

var DefaultConfigLocations = []string{
	"/etc/optnix/config.toml",
	// User config path filled in by init(), depending on `XDG_CONFIG_HOME` presence
	// optnix.toml in the current directory, if it exists
}

func init() {
	var homeDirPath string
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		homeDirPath = filepath.Join(xdgConfigHome, "optnix", "config.toml")
	} else if home := os.Getenv("HOME"); home != "" {
		homeDirPath = filepath.Join(home, ".config", "optnix", "config.toml")
	}

	if homeDirPath != "" {
		DefaultConfigLocations = append(DefaultConfigLocations, homeDirPath)
	}

	DefaultConfigLocations = append(DefaultConfigLocations, "optnix.toml")
}
