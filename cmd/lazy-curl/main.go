package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int
	items []string
	selected map[int]struct{}
}

func initialModel() model {
	return model {
		items: []string{"GET /api/users", "POST /api/users", "GET /api/settings"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "Enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Welcome to lazy-curl!\n\n"

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, item)
	}

	s += "\nPress j/k or up/down to move, q to quit.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There was an Error")
		os.Exit(1)
	}
}
