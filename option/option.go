package option

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"

	"github.com/google/shlex"
)

type NixosOption struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Type         string            `json:"type"`
	Default      *NixosOptionValue `json:"default"`
	Example      *NixosOptionValue `json:"example"`
	Location     []string          `json:"loc"`
	ReadOnly     bool              `json:"readOnly"`
	Declarations []string          `json:"declarations"`
}

type NixosOptionValue struct {
	Type string `json:"_type"`
	Text string `json:"text"`
}

type NixosOptionSource []NixosOption

func (o NixosOptionSource) String(i int) string {
	return o[i].Name
}

func (o NixosOptionSource) Len() int {
	return len(o)
}

func LoadOptions(r io.Reader) (NixosOptionSource, error) {
	var options []NixosOption

	d := json.NewDecoder(r)
	err := d.Decode(&options)
	if err != nil {
		return nil, err
	}

	return options, nil
}

var (
	elidedSetPattern  = regexp.MustCompile(`\{[ ]*\.\.\.[ ]*\}`)
	elidedListPattern = regexp.MustCompile(`\[[ ]*\.\.\.[ ]*\]`)
	angledPattern     = regexp.MustCompile(`«.*?»`)
)

func escapeNixEvalAnnotations(s string) string {
	s = elidedSetPattern.ReplaceAllString(s, `"«elided set»"`)
	s = elidedListPattern.ReplaceAllString(s, `"«elided list»"`)

	s = angledPattern.ReplaceAllStringFunc(s, func(s string) string {
		return fmt.Sprintf("%q", s)
	})

	return s
}

// Format a `nix-instantiate` or a `nix eval`-created string for pretty
// printing using a Nix code formatter command.
//
// The command passed must take the Nix code from `stdin` and pass it back
// out using `stdout`.
func FormatNixValue(formatterCmd string, evaluatedValue string) (string, error) {
	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString(escapeNixEvalAnnotations(evaluatedValue))

	argv, err := shlex.Split(formatterCmd)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(argv[0], argv[1:]...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = &stdin

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}
