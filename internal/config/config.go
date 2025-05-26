package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	MinScore     int64  `koanf:"min_score"`
	DebounceTime int64  `koanf:"debounce_time"`
	DefaultScope string `koanf:"default_scope"`

	Scopes map[string]Scope `koanf:"scopes"`
}

type Scope struct {
	Description     string `koanf:"description"`
	OptionsListFile string `koanf:"options-list-file"`
	OptionsListCmd  string `koanf:"options-list-cmd"`
	EvaluatorCmd    string `koanf:"evaluator"`
}

func NewConfig() *Config {
	return &Config{
		MinScore:     1,
		DebounceTime: 25,

		Scopes: make(map[string]Scope),
	}
}

func ParseConfig(location ...string) (*Config, error) {
	k := koanf.New(".")

	for _, loc := range location {
		if _, err := os.Stat(loc); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return nil, err
		}

		if err := k.Load(file.Provider(loc), toml.Parser()); err != nil {
			return nil, err
		}
	}

	cfg := NewConfig()

	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	for s, v := range c.Scopes {
		if v.OptionsListCmd == "" && v.OptionsListFile == "" {
			return fmt.Errorf("no option list source defined for scope %v", s)
		}
	}

	if c.DefaultScope != "" {
		foundScope := false
		for k := range c.Scopes {
			if k == c.DefaultScope {
				foundScope = true
				break
			}
		}
		if !foundScope {
			return fmt.Errorf("default config %v not found", c.DefaultScope)
		}
	}

	return nil
}

var DefaultConfigLocations = []string{
	"/etc/optnix/config.toml",
	// User config path filled in by init(), depending on `XDG_CONFIG_HOME` presence
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
}
