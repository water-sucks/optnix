package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/water-sucks/optnix/internal/config"
	"github.com/water-sucks/optnix/internal/logger"
)

func GenerateCompletions(cmd *cobra.Command, shell string) {
	switch shell {
	case "bash":
		_ = cmd.Root().GenBashCompletionV2(os.Stdout, true)
	case "zsh":
		_ = cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		_ = cmd.Root().GenFishCompletion(os.Stdout, true)
	}
}

func completeCompletionShells(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"bash", "fish", "zsh"}, cobra.ShellCompDirectiveDefault
}

func completeOptionsFromScope(scopeName *string) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 1 || *scopeName == "" {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		cfg := config.FromContext(cmd.Context())
		log := logger.FromContext(cmd.Context())

		scope, err := getScopeFromCfg(*scopeName, cfg)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		options, err := getOptionListFromScope(log, *scopeName, scope)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		names := make([]string, 0, len(options))
		for _, v := range options {
			if strings.HasPrefix(v.Name, toComplete) {
				names = append(names, v.Name)
			}
		}

		directive := cobra.ShellCompDirectiveNoFileComp
		if len(names) > 1 {
			directive |= cobra.ShellCompDirectiveNoSpace
		}

		return names, directive
	}
}

func completeScopes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg := config.FromContext(cmd.Context())

	scopes := []string{}
	for name := range cfg.Scopes {
		scopes = append(scopes, name)
	}

	return scopes, cobra.ShellCompDirectiveNoFileComp
}
