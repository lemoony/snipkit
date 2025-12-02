package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_ConfigDiffIntegration_FullFlow(t *testing.T) {
	// Real config data with meaningful differences
	oldConfig := `version: v1.0.0
config:
  editor: vim
  theme: default
  fuzzySearch: false`

	newConfig := `version: v1.1.0
config:
  editor: nvim
  theme: dracula
  fuzzySearch: true
  showPreview: true`

	// Render the confirmation message (tests full template pipeline)
	confirm := uimsg.ManagerConfigAddConfirm(oldConfig, newConfig)
	output := testutil.StripANSI(confirm.Header(style.NoopStyle, 120))

	// Verify diff is displayed
	assert.Contains(t, output, "BEFORE")
	assert.Contains(t, output, "AFTER")

	// Verify old values appear
	assert.Contains(t, output, "vim")
	assert.Contains(t, output, "default")
	assert.Contains(t, output, "false")

	// Verify new values appear
	assert.Contains(t, output, "nvim")
	assert.Contains(t, output, "dracula")
	assert.Contains(t, output, "true")
	assert.Contains(t, output, "showPreview")
}

func Test_ConfigDiffIntegration_MinimalConfig(t *testing.T) {
	oldConfig := "" // Empty for new installation
	newConfig := `version: v1.1.0
config:
  editor: nvim
  theme: dracula`

	confirm := uimsg.ConfigFileMigrationConfirm(oldConfig, newConfig)
	output := testutil.StripANSI(confirm.Header(style.NoopStyle, 120))

	// Should use "new config only" rendering
	assert.Contains(t, output, "NEW CONFIGURATION")
	assert.NotContains(t, output, "BEFORE")
	assert.NotContains(t, output, "AFTER")

	// Verify new config content is present
	assert.Contains(t, output, "nvim")
	assert.Contains(t, output, "dracula")
}

func Test_ConfigDiffIntegration_ContextCompression(t *testing.T) {
	// Build large old config (30 lines, mostly unchanged)
	oldConfig := buildLargeConfig(30, map[int]string{
		5:  "oldValue1",
		15: "oldValue2",
		25: "oldValue3",
	})

	// Build new config with only specific changes
	newConfig := buildLargeConfig(30, map[int]string{
		5:  "newValue1",
		15: "newValue2",
		25: "newValue3",
	})

	confirm := uimsg.ConfigFileMigrationConfirm(oldConfig, newConfig)
	output := testutil.StripANSI(confirm.Header(style.NoopStyle, 120))

	// Verify compression happened
	assert.Contains(t, output, "lines unchanged")

	// Verify changed lines are visible
	assert.Contains(t, output, "newValue1")
	assert.Contains(t, output, "newValue2")
	assert.Contains(t, output, "newValue3")

	// Verify not all 30 lines are in output (compression worked)
	outputLines := strings.Split(output, "\n")
	assert.Less(t, len(outputLines), 40) // Should be compressed from 30+ to much less
}

func Test_ConfigDiffIntegration_AllChangeTypes(t *testing.T) {
	oldConfig := `version: v1.0.0
config:
  removedField: value
  modifiedField: oldValue
  unchangedField: same`

	newConfig := `version: v1.0.0
config:
  addedField: newField
  modifiedField: newValue
  unchangedField: same`

	confirm := uimsg.ManagerConfigAddConfirm(oldConfig, newConfig)
	output := testutil.StripANSI(confirm.Header(style.NoopStyle, 120))

	// All types should be present
	assert.Contains(t, output, "removedField")   // Deletion
	assert.Contains(t, output, "addedField")     // Addition
	assert.Contains(t, output, "modifiedField")  // Modification
	assert.Contains(t, output, "unchangedField") // Context

	// Both old and new values for modified field
	assert.Contains(t, output, "oldValue")
	assert.Contains(t, output, "newValue")
}

// buildLargeConfig creates a config with specified number of lines and changes at specific line numbers.
func buildLargeConfig(lines int, changes map[int]string) string {
	var sb strings.Builder
	sb.WriteString("version: v1.0.0\n")
	sb.WriteString("config:\n")

	for i := 1; i <= lines; i++ {
		if value, ok := changes[i]; ok {
			sb.WriteString(fmt.Sprintf("  field%d: %s\n", i, value))
		} else {
			sb.WriteString(fmt.Sprintf("  field%d: defaultValue%d\n", i, i))
		}
	}

	return sb.String()
}
