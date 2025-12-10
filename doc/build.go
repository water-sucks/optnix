package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"snare.dev/optnix/option"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "build",
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
			HiddenDefaultCmd:  true,
		},
	}
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	moduleGenCmd := &cobra.Command{
		Use:   "gen-module-docs",
		Short: "Generate Markdown documentation for modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("generating module markdown documentation")

			generatedModulePath := filepath.Join("doc", "src", "usage", "generated-module.md")
			if err := generateModuleDocMarkdown(generatedModulePath); err != nil {
				return err
			}

			fmt.Println("generated module documentation for mdbook site")

			return nil
		},
	}

	rootCmd.AddCommand(moduleGenCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildModuleOptionsJSON() (option.NixosOptionSource, error) {
	buildModuleDocArgv := []string{"nix-build", "./doc/options.nix"}

	var buildModuleDocStdout bytes.Buffer

	buildModuleDocCmd := exec.Command(buildModuleDocArgv[0], buildModuleDocArgv[1:]...)
	buildModuleDocCmd.Stdout = &buildModuleDocStdout

	err := buildModuleDocCmd.Run()
	if err != nil {
		return nil, err
	}

	optionsDocFilename := strings.TrimSpace(buildModuleDocStdout.String())

	optionsDocFile, err := os.Open(optionsDocFilename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = optionsDocFile.Close() }()

	return option.LoadOptions(optionsDocFile)
}

func formatOptionMarkdown(opt option.NixosOption) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## `%s`\n\n", opt.Name))

	if opt.Description != "" {
		sb.WriteString(opt.Description + "\n\n")
	}

	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", opt.Type))

	if opt.Default != nil {
		sb.WriteString(fmt.Sprintf("**Default:** `%s`\n\n", opt.Default.Text))
	}

	if opt.Example != nil {
		sb.WriteString(fmt.Sprintf("**Example:** `%s`\n\n", opt.Example.Text))
	}

	return sb.String()
}

func generateModuleDocMarkdown(outputFilename string) error {
	options, err := buildModuleOptionsJSON()
	if err != nil {
		return err
	}

	var sb strings.Builder

	for _, opt := range options {
		sb.WriteString(formatOptionMarkdown(opt))
		sb.WriteString("\n")
	}

	err = os.WriteFile(outputFilename, []byte(sb.String()), 0o644)
	if err != nil {
		return err
	}

	return nil
}
