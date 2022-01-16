package app

import (
	"fmt"

	"github.com/lemoony/snipkit/internal/ui/picker"
)

func (a *appImpl) AddManager() {
	managerDescriptions := a.provider.ManagerDescriptions(a.config.Manager)
	listItems := make([]picker.Item, len(managerDescriptions))
	for i := range managerDescriptions {
		listItems[i] = picker.NewItem(managerDescriptions[i].Name, managerDescriptions[i].Description)
	}

	if index, ok := a.ui.ShowPicker(listItems); ok {
		fmt.Printf("Selected: %d\n", index)
	}
}
