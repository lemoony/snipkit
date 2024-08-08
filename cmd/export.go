package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lemoony/snipkit/internal/app"
)

var (
	exportFieldsFlag []string
	exportFieldMap   = map[string]app.ExportField{
		"id":         app.ExportFieldID,
		"title":      app.ExportFieldTitle,
		"content":    app.ExportFieldContent,
		"parameters": app.ExportFieldParameters,
	}
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports snippets on stdout",
	Long:  `Exports all snippets on stdout, optionally including parsed meta information like parameters etc..`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		fmt.Println(app.ExportSnippets(exportedFields()))
	},
}

func exportedFields() []app.ExportField {
	result := make([]app.ExportField, len(exportFieldsFlag))
	for i, v := range exportFieldsFlag {
		if exportField, ok := exportFieldMap[v]; ok {
			result[i] = exportField
		} else {
			panic("Unsupported field: " + v)
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.PersistentFlags().StringSliceVarP(
		&exportFieldsFlag,
		"fields",
		"f",
		[]string{"id", "title", "content", "parameters"},
		"Fields to be exported",
	)
}
