package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	buildOpts "github.com/water-sucks/optnix/internal/build"
	cmdUtils "github.com/water-sucks/optnix/internal/cmd/utils"
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

	// TODO: init logger

	// TODO: load configurations

	cmd := cobra.Command{
		Use:   "optnix [SCOPE] [OPTION]",
		Short: "optnix-cli",
		Long:  "optnix - a fast Nix module system options searcher",
		Args: func(cmd *cobra.Command, args []string) error {
			// Grab positional args [SCOPE] and [OPTION]
			argc := len(args)

			if argc == 0 {
				// TODO: maybe add additional scopes as a hint here?
				return cmdUtils.ArgParseError{Msg: "missing required argument [SCOPE]"}
			}
			opts.Scope = args[0]

			if !opts.Interactive && argc < 2 {
				return cmdUtils.ArgParseError{
					Msg:  "argument [NAME] is required for non-interactive mode",
					Hint: fmt.Sprintf(`try running "optnix %v 'option-name'"`, opts.Scope),
				}
			}

			opts.OptionInput = args[1]

			// Validation of flags
			if opts.JSON && opts.Interactive {
				return cmdUtils.ArgParseError{Msg: "--json and --interactive flags conflict"}
			}
			if opts.JSON && opts.ValueOnly {
				return cmdUtils.ArgParseError{Msg: "--json and --value-only flags conflict"}
			}
			if opts.ValueOnly && opts.Interactive {
				return cmdUtils.ArgParseError{Msg: "--interactive and --value-only flags conflict"}
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
			// TODO: load configurations passed on command line, and MinScore

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return CommandMain(cmd, &opts)
		},
	}

	// TODO: init with context

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.SetUsageTemplate(helpTemplate)

	boldRed := color.New(color.FgRed).Add(color.Bold)
	cmd.SetErrPrefix(boldRed.Sprint("error:"))

	cmd.Flags().BoolP("help", "h", false, "Show this help menu")
	cmd.Flags().Bool("version", false, "Display version information")

	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Show interactive search TUI for options")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "Output information in JSON format")
	cmd.Flags().Int64VarP(&opts.MinScore, "min-score", "m", 0, "Minimum score threshold for matching")
	cmd.Flags().StringSliceVarP(&opts.Config, "config", "c", nil, "Path to extra configuration `files` to load")
	cmd.Flags().BoolVarP(&opts.ValueOnly, "value-only", "v", false, "Only show option values")

	return &cmd
}

func CommandMain(cmd *cobra.Command, opts *CmdOptions) error {
	return nil
}

func Execute() {
	if err := MainCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
