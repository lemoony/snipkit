package picker

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2, 0, 2) //nolint:gomnd // magic number is okay for styling purposes
	delegate = list.NewDefaultDelegate()
)

type Item struct {
	title, desc string
}

func NewItem(title, desc string) Item {
	return Item{title: title, desc: desc}
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	choice *Item
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		} else if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = &i
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return docStyle.Render(m.list.View())
}

func ShowPicker(items []Item, options ...tea.ProgramOption) (int, bool) {
	listItems := make([]list.Item, len(items))
	for i := range items {
		listItems[i] = list.Item(items[i])
	}

	delegate.SetSpacing(1)

	m := model{list: list.New(listItems, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Which snippet manager should be added to your configuration"
	m.list.SetDelegate(delegate)
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)

	m.list.KeyMap.AcceptWhileFiltering.SetEnabled(true)
	m.list.KeyMap.AcceptWhileFiltering.SetHelp("↵", "apply")
	m.list.KeyMap.ShowFullHelp.SetEnabled(false)

	p := tea.NewProgram(&m, append(options, tea.WithAltScreen())...)

	if err := p.Start(); err != nil {
		panic(err)
	}

	if m.choice != nil {
		c := *m.choice
		for i := range m.list.Items() {
			if listItems[i] == c {
				return i, true
			}
		}
	}

	return -1, false
}
