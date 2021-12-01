package cli

import (
	"fmt"
	"os"

	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/utils"
)

func Info() error {
	system, err := utils.NewSystem()
	if err != nil {
		return err
	}

	snippetsLab, err := snippetslab.NewProvider(
		snippetslab.WithSystem(&system),
		snippetslab.WithTagsFilter([]string{"snipkit", "footag"}),
	)
	if err != nil {
		return err
	}

	allProviders := []providers.Provider{
		snippetsLab,
	}

	for _, provider := range allProviders {
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
