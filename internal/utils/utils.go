package utils

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/sahilm/fuzzy"
)

func FilterMinimumScoreMatches(matches []fuzzy.Match, minScore int64) []fuzzy.Match {
	for i, v := range matches {
		if v.Score < int(minScore) {
			return matches[:i]
		}
	}

	return matches
}

type ShellExecOutput struct {
	State  *os.ProcessState
	Stdout string
	Stderr string
}

func ExecShellAndCaptureOutput(commandStr string) (ShellExecOutput, error) {
	cmd := exec.Command("/bin/sh", "-c", commandStr)
	cmd.Env = os.Environ()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	result := ShellExecOutput{}

	err := cmd.Run()

	result.State = cmd.ProcessState
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	return result, err
}
