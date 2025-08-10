package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
	buildOpts "github.com/water-sucks/optnix/internal/build"
	cmdUtils "github.com/water-sucks/optnix/internal/cmd/utils"
	"github.com/water-sucks/optnix/internal/config"
	"github.com/water-sucks/optnix/internal/logger"
	"github.com/water-sucks/optnix/internal/utils"
	"github.com/water-sucks/optnix/option"
	"github.com/water-sucks/optnix/tui"
	"github.com/yarlson/pin"
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
	NonInteractive      bool
	Config              []string
	JSON                bool
	MinScore            int64
	ValueOnly           bool
	Scope               string
	ListScopes          bool
	GenerateCompletions string

	OptionInput string
}

func CreateCommand() *cobra.Command {
	opts := CmdOptions{}
	cmdCtx := context.Background()

	log := logger.NewLogger()
	cmdCtx = logger.WithLogger(cmdCtx, log)

	cmd := cobra.Command{
		Use:   "optnix -s [SCOPE] [OPTION-NAME]",
		Short: "optnix",
		Long:  "optnix - a fast Nix module system options searcher",
		Args: func(cmd *cobra.Command, args []string) error {
			argc := len(args)

			if opts.GenerateCompletions != "" {
				switch opts.GenerateCompletions {
				case "bash", "zsh", "fish":
				default:
					return cmdUtils.ErrorWithHint{
						Msg:  fmt.Sprintf("unsupported shell '%v'", opts.GenerateCompletions),
						Hint: "supported shells for completion are bash, zsh, or fish",
					}
				}
				return nil
			}

			if argc > 0 {
				opts.OptionInput = args[0]
			}

			// Imply `--non-interactive` for scripting output if not specified
			if opts.JSON || opts.ValueOnly {
				if cmd.Flags().Changed("non-interactive") && !opts.NonInteractive {
					return cmdUtils.ErrorWithHint{Msg: "--non-interactive is required when using output format flags"}
				}

				opts.NonInteractive = true
			}

			if opts.JSON && opts.ValueOnly {
				return cmdUtils.ErrorWithHint{Msg: "--json and --value-only flags conflict"}
			}

			if opts.NonInteractive && argc < 1 {
				scopeName := opts.Scope
				if scopeName == "" {
					scopeName = "[SCOPE]"
				}

				return cmdUtils.ErrorWithHint{
					Msg:  "argument [OPTION-NAME] is required for non-interactive mode",
					Hint: fmt.Sprintf(`try running "optnix -s %v [OPTION-NAME]"`, scopeName),
				}
			}

			return nil
		},
		Version:                    buildOpts.Version,
		SilenceUsage:               true,
		SuggestionsMinimumDistance: 1,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		ValidArgsFunction: completeOptionsFromScope(&opts.Scope),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.GenerateCompletions != "" {
				return nil
			}

			inCompletionMode := cmd.CalledAs() == cobra.ShellCompRequestCmd

			configLocations := append(config.DefaultConfigLocations, opts.Config...)

			cfg, err := config.ParseConfig(configLocations...)
			if err != nil {
				log.Errorf("failed to parse config: %v", err)
				return err
			}

			if !inCompletionMode {
				if err := cfg.Validate(); err != nil {
					return err
				}
			}

			if opts.Scope == "" {
				if cfg.DefaultScope == "" && !inCompletionMode && !opts.ListScopes {
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
			if err := commandMain(cmd, &opts); err != nil {
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
	cmd.Flags().BoolVarP(&opts.NonInteractive, "non-interactive", "n", false, "Do not show search TUI for options")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "Output information in JSON format")
	cmd.Flags().BoolVarP(&opts.ListScopes, "list-scopes", "l", false, "List available scopes and exit")
	cmd.Flags().Int64VarP(&opts.MinScore, "min-score", "m", 0, "Minimum `score` threshold for matching")
	cmd.Flags().StringSliceVarP(&opts.Config, "config", "c", nil, "Path to extra configuration `files` to load")
	cmd.Flags().BoolVarP(&opts.ValueOnly, "value-only", "v", false, "Only show option values")

	cmd.Flags().StringVar(&opts.GenerateCompletions, "completion", "", "Generate completions for a shell")
	_ = cmd.Flags().MarkHidden("completion")

	_ = cmd.RegisterFlagCompletionFunc("scope", completeScopes)
	_ = cmd.RegisterFlagCompletionFunc("completion", completeCompletionShells)

	return &cmd
}

var scopesListHeader = []string{"Name", "Description", "Origin"}

func centered(width int, s string) *bytes.Buffer {
	var b bytes.Buffer
	runeLen := utf8.RuneCountInString(s)

	if runeLen >= width {
		fmt.Fprint(&b, s)
		return &b
	}

	padding := width - runeLen
	left := padding / 2
	right := padding - left

	fmt.Fprintf(&b, "%s%s%s", strings.Repeat(" ", left), s, strings.Repeat(" ", right))
	return &b
}

func listScopes(cfg *config.Config) {
	const separator = " | "

	data := make([][]string, 0, len(cfg.Scopes))

	for s, v := range cfg.Scopes {
		origin := cfg.FieldOrigin(fmt.Sprintf("scopes.%v", s))
		data = append(data, []string{s, v.Description, origin})
	}

	colWidths := make([]int, len(scopesListHeader))

	for _, row := range data {
		for i, col := range row {
			if len(col) > colWidths[i] {
				colWidths[i] = len(col)
			}
		}
	}

	totalWidth := 0

	for i, col := range scopesListHeader {
		fmt.Print(centered(colWidths[i], col))
		totalWidth += colWidths[i]

		if i < len(scopesListHeader)-1 {
			fmt.Print(separator)
			totalWidth += len(separator)
		}
	}

	fmt.Println("\n" + strings.Repeat("-", totalWidth))

	for _, row := range data {
		for i, col := range row {
			format := fmt.Sprintf("%%-%ds", colWidths[i])
			fmt.Printf(format, col)
			if i < len(row)-1 {
				fmt.Print(separator)
			}
		}
		fmt.Println()
	}
}

func constructScopeFromConfig(scope *config.Scope, formatterCmd string) option.Scope {
	loader := func() (option.NixosOptionSource, error) {
		return scope.Load()
	}

	evaluator := constructEvaluatorFromScope(formatterCmd, scope)

	return option.Scope{
		Name:        scope.Name,
		Description: scope.Description,
		Loader:      loader,
		Evaluator:   evaluator,
	}
}

func constructEvaluatorFromScope(formatterCmd string, s *config.Scope) option.EvaluatorFunc {
	if s.EvaluatorCmd == "" {
		return nil
	}

	tmpl, err := template.New("eval").Parse(s.EvaluatorCmd)
	if err != nil {
		panic(fmt.Sprintf("evaluator should have been verified as valid at this point: %v", err))
	}

	return func(optionName string) (string, error) {
		var buf bytes.Buffer

		err := tmpl.Execute(&buf, map[string]string{
			"Option": optionName,
		})
		if err != nil {
			return "", err
		}

		cmdOutput, err := utils.ExecShellAndCaptureOutput(buf.String())
		if err != nil {
			return "", &option.AttributeEvaluationError{
				Attribute:        optionName,
				EvaluationOutput: strings.TrimSpace(cmdOutput.Stderr),
			}
		}

		output := cmdOutput.Stdout

		if formatterCmd != "" {
			if formatted, err := option.FormatNixValue(formatterCmd, output); err == nil {
				output = formatted
			}
		}

		value := strings.TrimSpace(output)

		return value, nil
	}
}

func commandMain(cmd *cobra.Command, opts *CmdOptions) error {
	if opts.GenerateCompletions != "" {
		GenerateCompletions(cmd, opts.GenerateCompletions)
		return nil
	}

	log := logger.FromContext(cmd.Context())
	cfg := config.FromContext(cmd.Context())

	if opts.ListScopes {
		listScopes(cfg)
		return nil
	}

	if !opts.NonInteractive {
		scopes := make([]option.Scope, 0, len(cfg.Scopes))
		for _, scope := range cfg.Scopes {
			actualScope := constructScopeFromConfig(&scope, cfg.FormatterCmd)
			scopes = append(scopes, actualScope)
		}

		return tui.OptionTUI(tui.OptionTUIArgs{
			Scopes:            scopes,
			SelectedScopeName: opts.Scope,
			MinScore:          cfg.MinScore,
			DebounceTime:      cfg.DebounceTime,
			InitialInput:      opts.OptionInput,
		})
	}

	var scope *option.Scope
	for _, s := range cfg.Scopes {
		if opts.Scope == s.Name {
			actualScope := constructScopeFromConfig(&s, cfg.FormatterCmd)
			scope = &actualScope
			break
		}
	}

	if scope == nil {
		err := fmt.Errorf("scope '%v' not found in configuration", opts.Scope)
		log.Errorf("%v", err)
		return err
	}

	spinner := pin.New("Loading...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorRed),
		pin.WithPosition(pin.PositionRight),
		pin.WithSpinnerFrames([]rune{'-', '\\', '|', '/'}),
		pin.WithWriter(os.Stderr),
	)
	cancelSpinner := spinner.Start(context.Background())
	defer cancelSpinner()

	spinner.UpdateMessage("Loading options...")

	options, err := scope.Loader()
	if err != nil {
		spinner.Stop()
		log.Errorf("%v", err)
		return err
	}

	spinner.UpdateMessage(fmt.Sprintf("Finding option %v...", opts.OptionInput))

	exactOptionMatchIdx := slices.IndexFunc(options, func(o option.NixosOption) bool {
		return o.Name == opts.OptionInput
	})
	if exactOptionMatchIdx != -1 {
		o := options[exactOptionMatchIdx]

		spinner.UpdateMessage("Evaluating option value...")
		var evaluatedValue string
		var evalErr error

		if scope.Evaluator != nil {
			evaluatedValue, evalErr = scope.Evaluator(o.Name)
		} else {
			evaluatedValue = "no evaluator configured for this scope"
		}

		spinner.Stop()

		if opts.JSON {
			displayOptionJson(&o, &evaluatedValue)
		} else if opts.ValueOnly {
			fmt.Printf("%v\n", evaluatedValue)
		} else {
			fmt.Print(o.PrettyPrint(&option.ValuePrinterInput{
				Value: evaluatedValue,
				Err:   evalErr,
			}))
		}

		return nil
	}

	spinner.Stop()

	msg := fmt.Sprintf("no exact match for query '%s' found", opts.OptionInput)
	err = fmt.Errorf("%v", msg)

	fuzzySearchResults := fuzzy.FindFrom(opts.OptionInput, options)
	if len(fuzzySearchResults) > 10 {
		fuzzySearchResults = fuzzySearchResults[:10]
	}

	fuzzySearchResults = utils.FilterMinimumScoreMatches(fuzzySearchResults, cfg.MinScore)

	if opts.JSON {
		displayErrorJson(msg, fuzzySearchResults)
		return err
	}

	log.Error(msg)
	if len(fuzzySearchResults) > 0 {
		log.Print("\nSome similar options were found:\n")
		for _, v := range fuzzySearchResults {
			log.Printf(" - %s\n", v.Str)
		}
	} else {
		log.Print("\nTry refining your search query.\n")
	}

	return err
}

type optionJsonOutput struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Value        *string  `json:"value"`
	Default      string   `json:"default"`
	Example      string   `json:"example"`
	Location     []string `json:"loc"`
	ReadOnly     bool     `json:"readOnly"`
	Declarations []string `json:"declarations"`
}

func displayOptionJson(o *option.NixosOption, evaluatedValue *string) {
	defaultText := ""
	if o.Default != nil {
		defaultText = o.Default.Text
	}

	exampleText := ""
	if o.Example != nil {
		exampleText = o.Example.Text
	}

	bytes, _ := json.MarshalIndent(optionJsonOutput{
		Name:         o.Name,
		Description:  o.Description,
		Type:         o.Type,
		Value:        evaluatedValue,
		Default:      defaultText,
		Example:      exampleText,
		Location:     o.Location,
		ReadOnly:     o.ReadOnly,
		Declarations: o.Declarations,
	}, "", "  ")
	fmt.Printf("%v\n", string(bytes))
}

type errorJsonOutput struct {
	Message        string   `json:"message"`
	SimilarOptions []string `json:"similar_options"`
}

func displayErrorJson(msg string, matches fuzzy.Matches) {
	matchedStrings := make([]string, len(matches))
	for i, match := range matches {
		matchedStrings[i] = match.Str
	}

	bytes, _ := json.MarshalIndent(errorJsonOutput{
		Message:        msg,
		SimilarOptions: matchedStrings,
	}, "", "  ")
	fmt.Printf("%v\n", string(bytes))
}

func Execute() {
	if err := CreateCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
