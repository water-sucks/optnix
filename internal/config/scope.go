package config

import (
	"encoding/json"
	"fmt"
	"os"

	"snare.dev/optnix/internal/utils"
	"snare.dev/optnix/option"
)

type Scope struct {
	Name            string `koanf:"-"`
	Description     string `koanf:"description"`
	OptionsListFile string `koanf:"options-list-file"`
	OptionsListCmd  string `koanf:"options-list-cmd"`
	EvaluatorCmd    string `koanf:"evaluator"`
}

func (s Scope) Load() (option.NixosOptionSource, error) {
	if s.OptionsListFile != "" {
		optionsFile, err := os.Open(s.OptionsListFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open options file: %v", err)
		} else {
			defer func() { _ = optionsFile.Close() }()

			l, err := option.LoadOptions(optionsFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load options using file strategy: %v", err)
			} else {
				return l, nil
			}
		}
	}

	if s.OptionsListCmd != "" {
		l, err := runGenerateOptionListCmd(s.OptionsListCmd)
		if err != nil {
			return nil, fmt.Errorf("failed to run options cmd: %v", err)
		}

		return l, nil
	}

	return nil, fmt.Errorf("no options found through all strategies for scope '%v'", s.Name)
}

func runGenerateOptionListCmd(commandStr string) (option.NixosOptionSource, error) {
	cmdOutput, err := utils.ExecShellAndCaptureOutput(commandStr)
	if err != nil {
		return nil, err
	}

	var l option.NixosOptionSource
	err = json.Unmarshal([]byte(cmdOutput.Stdout), &l)
	if err != nil {
		return nil, err
	}

	return l, nil
}
