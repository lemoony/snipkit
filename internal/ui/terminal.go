package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/kballard/go-shellquote"

	"github.com/lemoony/snippet-kit/internal/model"
)

const (
	envEditor            = "EDITOR"
	envVisual            = "VISUAL"
	defaultEditor        = "vim"
	defaultEditorWindows = "notepad"
	windows              = "windows"
)

type Terminal interface {
	PrintMessage(message string)
	PrintError(message string)
	Confirm(message string) (bool, error)
	OpenEditor(path string, preferredEditor string) error
	ShowLookup(snippets []model.Snippet) (int, error)
	ShowParameterForm(parameters []model.Parameter) ([]string, error)
}

type ActualCLI struct {
	stdio terminal.Stdio
}

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
	if err := survey.AskOne(prompt, &confirmed, survey.WithStdio(a.stdio.In, a.stdio.Out, a.stdio.Err)); err != nil {
		return false, err
	}
	return confirmed, nil
}

func (a ActualCLI) OpenEditor(path string, preferredEditor string) error {
	editor := getEditor(preferredEditor)

	fmt.Println("---> " + editor)

	args, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	args = append(args, path)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // subprocess launched with a potential tainted input
	cmd.Stdin = a.stdio.In
	cmd.Stdout = a.stdio.Out
	cmd.Stderr = a.stdio.Err

	err = cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (a ActualCLI) ShowLookup(snippets []model.Snippet) (int, error) {
	return showLookup(snippets)
}

func (a ActualCLI) ShowParameterForm(parameters []model.Parameter) ([]string, error) {
	return showParameterForm(parameters)
}

func getEditor(preferred string) string {
	result := defaultEditor
	if runtime.GOOS == windows {
		result = defaultEditorWindows
	}

	preferred = strings.TrimSpace(preferred)
	if preferred != "" {
		result = preferred
	} else if v := os.Getenv(envVisual); v != "" {
		result = v
	} else if e := os.Getenv(envEditor); e != "" {
		result = e
	}

	return result
}
