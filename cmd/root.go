package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/google/shlex"
	"github.com/spf13/cobra"
	buildOpts "github.com/water-sucks/optnix/internal/build"
	cmdUtils "github.com/water-sucks/optnix/internal/cmd/utils"
	"github.com/water-sucks/optnix/internal/config"
	"github.com/water-sucks/optnix/internal/logger"
	"github.com/water-sucks/optnix/option"
)

const helpTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{if not .AllChildCommandsHaveGroup}}

Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{range $group := .Groups}}

{{.Title}}:{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

type CmdOptions struct {
	Interactive bool
	Config      []string
	JSON        bool
	MinScore    int64
	ValueOnly   bool
	Scope       string

	OptionInput string
}

func MainCommand() *cobra.Command {
	opts := CmdOptions{}
	cmdCtx := context.Background()

	log := logger.NewLogger()
	cmdCtx = logger.WithLogger(cmdCtx, log)

	cmd := cobra.Command{
		Use:   "optnix -s [SCOPE] [OPTION-NAME]",
		Short: "optnix-cli",
		Long:  "optnix - a fast Nix module system options searcher",
		Args: func(cmd *cobra.Command, args []string) error {
			argc := len(args)

			if !opts.Interactive && argc < 1 {
				return cmdUtils.ErrorWithHint{
					Msg:  "argument [OPTION-NAME] is required for non-interactive mode",
					Hint: fmt.Sprintf(`try running "optnix -s %v [OPTION-NAME]"`, opts.Scope),
				}
			}

			opts.OptionInput = args[0]

			// Validation of flags
			if opts.JSON && opts.Interactive {
				return cmdUtils.ErrorWithHint{Msg: "--json and --interactive flags conflict"}
			}
			if opts.JSON && opts.ValueOnly {
				return cmdUtils.ErrorWithHint{Msg: "--json and --value-only flags conflict"}
			}
			if opts.ValueOnly && opts.Interactive {
				return cmdUtils.ErrorWithHint{Msg: "--interactive and --value-only flags conflict"}
			}

			return nil
		},
		Version:                    buildOpts.Version,
		SilenceUsage:               true,
		SuggestionsMinimumDistance: 1,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configLocations := append(config.DefaultConfigLocations, opts.Config...)

			cfg, err := config.ParseConfig(configLocations...)
			if err != nil {
				log.Errorf("failed to parse config: %v", err)
				return err
			}

			if opts.Scope == "" {
				if cfg.DefaultScope == "" {
					return cmdUtils.ErrorWithHint{
						Msg:  "no scope was provided and no default scope is set in the configuration",
						Hint: "either set a default configuration or specify one with -s",
					}
				}

				opts.Scope = cfg.DefaultScope
			}

			if opts.MinScore != 0 {
				cfg.MinScore = opts.MinScore
			}

			cmdCtx = config.WithConfig(cmd.Context(), cfg)
			cmd.SetContext(cmdCtx)

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := CommandMain(cmd, &opts); err != nil {
				os.Exit(1)
			}
		},
	}

	cmd.SetContext(cmdCtx)

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.SetUsageTemplate(helpTemplate)

	boldRed := color.New(color.FgRed).Add(color.Bold)
	cmd.SetErrPrefix(boldRed.Sprint("error:"))

	cmd.Flags().BoolP("help", "h", false, "Show this help menu")
	cmd.Flags().Bool("version", false, "Display version information")

	cmd.Flags().StringVarP(&opts.Scope, "scope", "s", "", "Scope `name` to use (required)")
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Show interactive search TUI for options")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "Output information in JSON format")
	cmd.Flags().Int64VarP(&opts.MinScore, "min-score", "m", 0, "Minimum `score` threshold for matching")
	cmd.Flags().StringSliceVarP(&opts.Config, "config", "c", nil, "Path to extra configuration `files` to load")
	cmd.Flags().BoolVarP(&opts.ValueOnly, "value-only", "v", false, "Only show option values")

	return &cmd
}

func runGenerateOptionListCmd(commandStr string) (option.NixosOptionSource, error) {
	argv, err := shlex.Split(commandStr)
	if err != nil {
		return nil, fmt.Errorf("malformed command: %w", err)
	}

	var stdout bytes.Buffer

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdout = &stdout

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	var l option.NixosOptionSource
	err = json.Unmarshal(stdout.Bytes(), &l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func CommandMain(cmd *cobra.Command, opts *CmdOptions) error {
	log := logger.FromContext(cmd.Context())
	cfg := config.FromContext(cmd.Context())

	var scope *config.Scope
	for s, v := range cfg.Scopes {
		if s == opts.Scope {
			scope = &v
			break
		}
	}

	if scope == nil {
		err := fmt.Errorf("scope '%v' not found in configuration", opts.Scope)
		log.Errorf("%v", err)
		return err
	}

	var optionsList option.NixosOptionSource

	if scope.OptionsListFile != "" {
		optionsFile, err := os.Open(scope.OptionsListFile)
		if err != nil {
			log.Errorf("failed to open options file: %v", err)
		} else {
			defer func() { _ = optionsFile.Close() }()

			l, err := option.LoadOptions(optionsFile)
			if err != nil {
				log.Errorf("failed to load options using file strategy: %v", err)
				log.Info("attempting to load using command strategy instead")
			} else {
				optionsList = l
			}
		}
	}

	if len(optionsList) == 0 && scope.OptionsListCmd != "" {
		l, err := runGenerateOptionListCmd(scope.OptionsListCmd)
		if err != nil {
			log.Errorf("failed to run options cmd: %v", err)
			return err
		}

		optionsList = l
	}

	if optionsList == nil {
		err := fmt.Errorf("no options found through all strategies for scope '%v'", opts.Scope)
		log.Errorf("%v", err)
		return err
	}

	log.Infof("options list length: %v", len(optionsList))

	return nil
}

func Execute() {
	if err := MainCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
