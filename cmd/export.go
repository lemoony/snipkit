package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lemoony/snipkit/internal/app"
)

var (
	exportFormatFlag string
	exportFormatMap  = map[string]app.ExportFormat{
		"json":        app.ExportFormatJSON,
		"json-pretty": app.ExportFormatPrettyJSON,
		"xml":         app.ExportFormatXML,
	}

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
	Long:  `Exports all snippets on stdout as JSON including parsed meta information like parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		fmt.Println(app.ExportSnippets(exportedFields(), exportFormat()))
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

func exportFormat() app.ExportFormat {
	if format, ok := exportFormatMap[exportFormatFlag]; ok {
		return format
	}
	panic("Unsupported export format: " + exportFormatFlag)
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

	exportCmd.PersistentFlags().StringVarP(
		&exportFormatFlag,
		"output",
		"o",
		"json",
		"Output format. One of: json,json-pretty,xml",
	)
}
