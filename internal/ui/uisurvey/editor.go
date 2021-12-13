package uisurvey

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kballard/go-shellquote"
)

const (
	defaultEditor        = "vim"
	defaultEditorWindows = "notepad"
)

func Edit(path string, preferredEditor string) error {
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
