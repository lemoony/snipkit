package main

import (
	"fmt"
	"strings"

	internalModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/form"
)

func main() {
	fields := []internalModel.Parameter{
		{Key: "message", Name: "Message", Description: "What to print first"},
		{
			Key:         "application",
			Description: "A second information for the terminal",
			Values: []string{
				"The Romans learned from the Greeks",
				"probably marmelada",
				"by the French name cotignac",
				"option 4",
				"optopn 5",
				"option 6",
				"option 7",
				"option 8",
			},
		},
		{
			Key:          "Statement",
			Description:  "A description",
			DefaultValue: "default value",
		},
	}

	if values, ok := form.Show(fields, "OK"); ok {
		fmt.Printf("Apply: %s\n", strings.Join(values, ","))
	} else {
		fmt.Println("Dont apply")
	}
}
