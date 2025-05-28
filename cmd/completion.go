package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
