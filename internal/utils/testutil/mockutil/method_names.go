package mockutil

import (
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/utils/sliceutil"
)

const (
	Print        = "Print"
	PrintMessage = "PrintMessage"
	PrintError   = "PrintError"

	ApplyConfig = "ApplyConfig"

	Confirmation = "Confirmation"

	ShowPicker                 = "ShowPicker"
	ShowAssistantPrompt        = "ShowAssistantPrompt"
	ShowAssistantScriptPreview = "ShowAssistantScriptPreview"
	ShowAssistantWizard        = "ShowAssistantWizard"
	ShowSpinner                = "ShowSpinner"
	ShowParameterForm          = "ShowParameterForm"
	OpenEditor                 = "OpenEditor"

	Query                = "Query"
	ValidateConfig       = "Initialize"
	SaveAssistantSnippet = "SaveAssistantSnippet"
)

func FindMethodCall(method string, calls []mock.Call) *mock.Call {
	if call, ok := sliceutil.FindElement(calls, func(call mock.Call) bool {
		return call.Method == method
	}); ok {
		return &call
	}
	panic("Failed to find method call for " + method)
}
