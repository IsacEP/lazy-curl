package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int
	items []string
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
		s += fmt.Sprintf("%s %s \n", cursor, item)
	}

	s += "\nPress j/k or up/down to move, q to quit.\n"
	return s
}

func main() {
	initialState := model {
		items: []string{"GET /api/users", "POST /api/users", "GET /api/settings"},
	}

	p := tea.NewProgram(initialState)

	if _, err := p.Run(); err != nil {
		fmt.Printf("There was an Error")
		os.Exit(1)
	}
}
