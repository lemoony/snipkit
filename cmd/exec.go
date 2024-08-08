package cmd

import (
	"regexp"

	"github.com/spf13/cobra"

	"github.com/lemoony/snipkit/internal/model"
)

var (
	execCmdPrintFlag      = false
	execCmdConfirmFlag    = false
	execCmdIDFlag         string
	execCmdParametersFlag []string

	parameterValueRegex = regexp.MustCompile(`^(?P<key>[a-zA-Z_][a-zA-Z0-9_]*)=(?P<value>.*)$`)
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet directly from the terminal",
	Long:  `Execute a snippet directly from the terminal. The output of the commands will be visibile in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())

		if execCmdIDFlag == "" {
			app.LookupAndExecuteSnippet(execCmdConfirmFlag, execCmdPrintFlag)
		} else {
			app.FindScriptAndExecuteWithParameters(execCmdIDFlag, toParameterValues(execCmdParametersFlag), execCmdConfirmFlag, execCmdPrintFlag)
		}
	},
}

func toParameterValues(flagValues []string) []model.ParameterValue {
	result := make([]model.ParameterValue, len(flagValues))
	for i, v := range flagValues {
		match := parameterValueRegex.FindAllStringSubmatch(v, -1)
		if match == nil || len(match) != 1 {
			panic("Invalid parameter value: " + v)
		}
		result[i] = model.ParameterValue{Key: match[0][1], Value: match[0][2]}
	}
	return result
}

func init() {
	execCmd.PersistentFlags().BoolVar(
		&execCmdPrintFlag,
		"print",
		false,
		"print the command before execution on stdout",
	)

	execCmd.PersistentFlags().BoolVar(
		&execCmdConfirmFlag,
		"confirm",
		false,
		"the command is printed on stdout before execution for confirmation",
	)

	execCmd.PersistentFlags().StringVar(
		&execCmdIDFlag,
		"id",
		"",
		"ID of the snippet to execute",
	)

	execCmd.PersistentFlags().StringArrayVarP(
		&execCmdParametersFlag,
		"param",
		"p",
		[]string{},
		"Parameter values to be passed to the snippet",
	)

	rootCmd.AddCommand(execCmd)
}
