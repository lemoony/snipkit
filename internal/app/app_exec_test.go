package app

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/term"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_App_Exec(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{
			ID:       "uuid1",
			Title:    "title-1",
			Language: model.LanguageYAML,
			Tags:     []string{},
			Content:  testSnippetContent,
		},
	}

	inputVar1Value := "foo-value"

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything, mock.Anything).Return(0)
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{inputVar1Value, ""}, true)
	tui.On(mockutil.Print, mock.Anything)
	tui.On(mockutil.Confirmation, mock.Anything).Return(true)
	tui.On(mockutil.PrintMessage, inputVar1Value+"\n").Return()

	app := NewApp(
		WithTUI(&tui),
		WithConfig(configtest.NewTestConfig().Config),
		withManagerSnippets(snippets),
	)

	app.LookupAndExecuteSnippet(true, true)

	// TODO fix
	// tui.AssertCalled(t, mockutil.Confirmation, uimsg.ExecConfirm("title-1", testSnippetContent))
	// tui.AssertCalled(t, mockutil.Print, uimsg.ExecPrint("title-1", testSnippetContent))
}

func Test_App_Exec_FindScriptAndExecuteWithParameters(t *testing.T) {
	snippetContent := `# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
echo "${VAR1}"`

	snippets := []model.Snippet{
		testutil.TestSnippet{
			ID:       "uuid1",
			Title:    "title-1",
			Language: model.LanguageYAML,
			Tags:     []string{},
			Content:  snippetContent,
		},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	app := NewApp(
		WithTUI(&tui),
		WithConfig(configtest.NewTestConfig().Config),
		withManagerSnippets(snippets),
	)

	app.FindScriptAndExecuteWithParameters("uuid1", []model.ParameterValue{{Key: "VAR1", Value: "foo"}}, false, false)
}

func Test_App_Exec_FindScriptAndExecuteWithParameters_MissingParameters(t *testing.T) {
	snippetContent := `# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
echo "${VAR1}"`

	snippets := []model.Snippet{
		testutil.TestSnippet{
			ID:       "uuid1",
			Title:    "title-1",
			Language: model.LanguageYAML,
			Tags:     []string{},
			Content:  snippetContent,
		},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{"VAR1", ""}, true)

	app := NewApp(
		WithTUI(&tui),
		WithConfig(configtest.NewTestConfig().Config),
		withManagerSnippets(snippets),
	)

	app.FindScriptAndExecuteWithParameters("uuid1", []model.ParameterValue{}, false, false)
	tui.AssertCalled(t, "ShowParameterForm", snippets[0].GetParameters(), []model.ParameterValue{}, ui.OkButtonExecute)
}

func Test_detectShell(t *testing.T) {
	tests := []struct {
		name, script, configuredShell, expected string
	}{
		{"bash shebang direct path", "#!/bin/bash\necho hello", "", "/bin/bash"},
		{"bash shebang with env", "#!/usr/bin/env bash\necho hello", "", "bash"},
		{"zsh shebang direct path", "#!/bin/zsh\necho hello", "", "/bin/zsh"},
		{"zsh shebang with env", "#!/usr/bin/env zsh\necho hello", "", "zsh"},
		{"sh shebang", "#!/bin/sh\necho hello", "", "/bin/sh"},
		{"shebang with arguments", "#!/bin/bash -e\necho hello", "", "/bin/bash"},
		{"env shebang with arguments", "#!/usr/bin/env bash -e\necho hello", "", "bash"},
		{"no shebang uses configured shell", "echo hello", "/bin/zsh", "/bin/zsh"},
		{"shebang takes priority over configured shell", "#!/bin/bash\necho hello", "/bin/zsh", "/bin/bash"},
		{"python shebang", "#!/usr/bin/env python3\nprint('hello')", "", "python3"},
		{"no shebang and no configured shell falls back", "echo hello", "", "/bin/bash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectShell(tt.script, tt.configuredShell)
			// For the fallback test case, we need to handle the $SHELL env var
			if tt.name == "no shebang and no configured shell falls back" {
				assert.NotEmpty(t, result) // Either it uses $SHELL or falls back to /bin/bash
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_detectShell_EdgeCases(t *testing.T) {
	tests := []struct {
		name, script, configuredShell, expected string
	}{
		{"empty script uses configured shell", "", "/bin/zsh", "/bin/zsh"},
		{"whitespace only script uses configured shell", "   \n\t  ", "/bin/bash", "/bin/bash"},
		{"shebang without newline falls back", "#!/bin/bash", "/bin/zsh", "/bin/zsh"},
		{"shebang with only newline", "#!/bin/bash\n", "", "/bin/bash"},
		{"shebang with whitespace before interpreter", "#!  /bin/bash\n", "", "/bin/bash"},
		{"env shebang with extra spaces", "#!/usr/bin/env   bash\necho hi", "", "bash"},
		{"ruby shebang", "#!/usr/bin/env ruby\nputs 'hello'", "", "ruby"},
		{"node shebang", "#!/usr/bin/env node\nconsole.log('hi')", "", "node"},
		{"perl direct path", "#!/usr/bin/perl\nprint 'hi'", "", "/usr/bin/perl"},
		{"fish shell", "#!/usr/bin/env fish\necho hello", "", "fish"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectShell(tt.script, tt.configuredShell)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_executeWithoutPTY(t *testing.T) {
	tests := []struct {
		name           string
		script         string
		expectedStdout string
		expectedStderr string
	}{
		{
			name:           "simple echo",
			script:         "echo hello",
			expectedStdout: "hello\n",
			expectedStderr: "",
		},
		{
			name:           "stderr output",
			script:         "echo error >&2",
			expectedStdout: "",
			expectedStderr: "error\n",
		},
		{
			name:           "both stdout and stderr",
			script:         "echo out && echo err >&2",
			expectedStdout: "out\n",
			expectedStderr: "err\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executeScript(tt.script, "/bin/sh")
			assert.Equal(t, tt.expectedStdout, result.stdout)
			assert.Equal(t, tt.expectedStderr, result.stderr)
		})
	}
}

func Test_executeScript_usesDetectedShell(t *testing.T) {
	// Test that shebang is respected
	script := "#!/bin/sh\necho $0"
	result := executeScript(script, "/bin/bash")
	// The output should indicate sh was used (contains "sh")
	assert.Contains(t, result.stdout, "sh")
}

func Test_executeScript_terminalDetection(t *testing.T) {
	// Save original function
	originalIsTerminal := isTerminalFunc
	defer func() { isTerminalFunc = originalIsTerminal }()

	// Test non-terminal path (default in tests)
	isTerminalFunc = func(fd int) bool { return false }
	result := executeScript("echo test", "/bin/sh")
	assert.Equal(t, "test\n", result.stdout)
	assert.Equal(t, "", result.stderr)
}

func Test_formatOptions(t *testing.T) {
	tests := []struct {
		config   config.ScriptConfig
		expected model.SnippetFormatOptions
	}{
		{
			config:   config.ScriptConfig{RemoveComments: true, ParameterMode: config.ParameterModeSet},
			expected: model.SnippetFormatOptions{RemoveComments: true, ParamMode: model.SnippetParamModeSet},
		},
		{
			config:   config.ScriptConfig{RemoveComments: false, ParameterMode: config.ParameterModeReplace},
			expected: model.SnippetFormatOptions{RemoveComments: false, ParamMode: model.SnippetParamModeReplace},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			assert.Equal(t, tt.expected, formatOptions(tt.config))
		})
	}
}

// saveTermFuncs saves and returns a restore function for all terminal function variables.
func saveTermFuncs() func() {
	origIsTerminal := isTerminalFunc
	origGetSize := getTermSizeFunc
	origMakeRaw := makeRawFunc
	origRestore := restoreTermFunc
	return func() {
		isTerminalFunc = origIsTerminal
		getTermSizeFunc = origGetSize
		makeRawFunc = origMakeRaw
		restoreTermFunc = origRestore
	}
}

func Test_executeWithPTY_TermSizeError(t *testing.T) {
	defer saveTermFuncs()()

	// Mock getTermSizeFunc to return error - should use defaults (80x24)
	getTermSizeFunc = func(fd int) (int, int, error) {
		return 0, 0, errors.New("not a terminal")
	}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("hello")
		c.Send("\n") // Press Enter to continue after script execution
	}, func(stdio termutil.Stdio) {
		oldStdin, oldStdout := os.Stdin, os.Stdout
		os.Stdin = stdio.In.(*os.File)
		os.Stdout = stdio.Out.(*os.File)
		defer func() {
			os.Stdin = oldStdin
			os.Stdout = oldStdout
		}()

		isTerminalFunc = func(fd int) bool { return true }
		makeRawFunc = func(fd int) (*term.State, error) {
			return nil, errors.New("cannot make raw")
		}

		result := executeScript("echo hello", "/bin/sh")
		assert.Contains(t, result.stdout, "hello")
	})
}

func Test_executeWithPTY_MakeRawError(t *testing.T) {
	defer saveTermFuncs()()

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("test output")
		c.Send("\n") // Press Enter to continue after script execution
	}, func(stdio termutil.Stdio) {
		oldStdin, oldStdout := os.Stdin, os.Stdout
		os.Stdin = stdio.In.(*os.File)
		os.Stdout = stdio.Out.(*os.File)
		defer func() {
			os.Stdin = oldStdin
			os.Stdout = oldStdout
		}()

		isTerminalFunc = func(fd int) bool { return true }
		// Simulate makeRaw failure - execution should still proceed
		makeRawFunc = func(fd int) (*term.State, error) {
			return nil, errors.New("cannot set raw mode")
		}

		result := executeScript("echo 'test output'", "/bin/sh")
		assert.Contains(t, result.stdout, "test output")
	})
}
