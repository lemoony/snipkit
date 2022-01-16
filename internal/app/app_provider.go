package app

import (
	"fmt"

	"github.com/lemoony/snipkit/internal/ui/picker"
)

func (a *appImpl) AddProvider() {
	providerDescriptions := a.providersBuilder.ProviderDescriptions(a.config.Providers)
	listItems := make([]picker.Item, len(providerDescriptions))
	for i := range providerDescriptions {
		listItems[i] = picker.NewItem(providerDescriptions[i].Name, providerDescriptions[i].Description)
	}

	if index, ok := a.ui.ShowPicker(listItems); ok {
		fmt.Printf("Selected: %d\n", index)
	}
}
