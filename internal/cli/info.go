package cli

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	app "github.com/lemoony/snippet-kit/internal/app"
)

func Info(v *viper.Viper) error {
	snipkit, err := app.NewApp(v)
	if snipkit == nil || err != nil {
		return err
	}

	for _, provider := range snipkit.Providers {
		for _, line := range provider.Info().Lines {
			if line.IsError {
				_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", line.Key, line.Value)
			} else {
				fmt.Printf("%s: %s\n", line.Key, line.Value)
			}
		}
	}
	return nil
}
