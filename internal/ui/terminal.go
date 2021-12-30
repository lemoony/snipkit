package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kballard/go-shellquote"

	"github.com/lemoony/snippet-kit/internal/model"
)

const (
	defaultEditor        = "vim"
	defaultEditorWindows = "notepad"
)

type Terminal interface {
	PrintMessage(message string)
	PrintError(message string)
	Confirm(message string) (bool, error)
	OpenEditor(path string, preferredEditor string) error
	ShowLookup(snippets []model.Snippet) (int, error)
	ShowParameterForm(parameters []model.Parameter) ([]string, error)
}

type ActualCLI struct{}

func NewTerminal() Terminal {
	return ActualCLI{}
}

func (a ActualCLI) PrintMessage(msg string) {
	fmt.Printf(msg + "\n")
}

func (a ActualCLI) PrintError(msg string) {
	fmt.Printf(msg + "\n")
}

func (a ActualCLI) Confirm(message string) (bool, error) {
	confirmed := false
	prompt := &survey.Confirm{Message: message}
	if err := survey.AskOne(prompt, &confirmed); err != nil {
		return false, err
	}
	return confirmed, nil
}

func (a ActualCLI) OpenEditor(path string, preferredEditor string) error {
	editor := defaultEditor
	if runtime.GOOS == "windows" {
		editor = defaultEditorWindows
	}

	preferredEditor = strings.TrimSpace(preferredEditor)
	if preferredEditor != "" {
		editor = preferredEditor
	} else if v := os.Getenv("VISUAL"); v != "" {
		editor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}

	args, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	args = append(args, path)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // subprocess launched with a potential tainted input
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (a ActualCLI) ShowLookup(snippets []model.Snippet) (int, error) {
	return showLookup(snippets)
}

func (a ActualCLI) ShowParameterForm(parameters []model.Parameter) ([]string, error) {
	return showParameterForm(parameters)
}
