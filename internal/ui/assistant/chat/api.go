package chat

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lemoony/snipkit/internal/ui/style"
)

// ShowUnifiedChat shows the unified chat interface that handles prompt input,
// script generation, action selection, and execution in a single view.
//
// The unified chat supports multiple modes:
// - Input mode: User types a new prompt
// - Generating mode: Script is being generated asynchronously
// - Action menu mode: Script is ready, user chooses an action
// - Post-execution mode: Script has been executed, user chooses next step
//
// Returns:
//   - script: The generated script (interface{})
//   - parameterValues: Collected parameter values ([]string)
//   - action: The user's chosen action (PreviewAction)
//   - latestPrompt: The user's entered prompt (string) - set when action is PreviewActionRevise
//   - saveFilename: Filename for saving (string)
//   - saveSnippetName: Snippet name for saving (string)
func ShowUnifiedChat(
	config UnifiedConfig,
	styler style.Style,
	teaOptions ...tea.ProgramOption,
) (interface{}, []string, PreviewAction, string, string, string) {
	m := newUnifiedChatModel(config, styler)

	if teaModel, err := tea.NewProgram(m, append(teaOptions, tea.WithAltScreen(), tea.WithMouseCellMotion())...).Run(); err != nil {
		return nil, nil, PreviewActionCancel, "", "", ""
	} else if resultModel, ok := teaModel.(*unifiedChatModel); ok {
		return resultModel.generatedScript,
			resultModel.parameterValues,
			resultModel.action,
			resultModel.latestPrompt,
			resultModel.saveFilename,
			resultModel.saveSnippetName
	}

	return nil, nil, PreviewActionCancel, "", "", ""
}
