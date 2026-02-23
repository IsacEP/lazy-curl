package ui

import (
	"fmt"

	"github.com/IsacEP/lazy-curl/internal/client"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).MarginBottom(1)

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)

	itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0"))

	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7CCB")).Italic(true)

	activeBorder = lipgloss.Color("#ff5ec9")
	inactiveBorder = lipgloss.Color("#91579b")

	baseBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

type Model struct {
	cursor int
	items []string
	selected map[int]struct{}
	response string
	textInput textinput.Model
	isTyping  bool
	focusedPane int
	responseBody string
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Some Random Text Here"
	ti.CharLimit = 156
	ti.Width = 50

	return Model {
		items:     []string{"TEST", "GET /api/users", "POST /api/users", "GET /api/settings"},
		selected:  make(map[int]struct{}),
		response:  "Ready to use server",
		textInput: ti,
		isTyping:  false,
		focusedPane: 0,
		responseBody: "response body",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case client.ServerStatusMsg:
		if msg.Err != nil {
			m.response = fmt.Sprintf("Error checking %s: %v", msg.URL, msg.Err)
			m.responseBody = msg.Err.Error()
		} else {
			m.response = fmt.Sprintf("Success! %s returned status %d", msg.URL, msg.Status)
			m.responseBody = "Status Code: " + fmt.Sprint(msg.Status) + "\n"
			m.responseBody += "JSON: " + fmt.Sprint(msg.Body) 
		}
		return m, nil

	case tea.KeyMsg:
		if m.isTyping {
			switch msg.String() {
			case "enter":
				val := m.textInput.Value()
				if val != "" {
					m.items = append(m.items, val)
					m.textInput.Reset()
				}
				m.isTyping = false
				m.textInput.Blur()
				return m, nil

			case "esc":
				m.isTyping = false
				m.textInput.Blur()
				return m, nil
			}

			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.focusedPane == 0 && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.focusedPane == 0 && m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case " ":
			if m.focusedPane == 0 {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}

		case "enter":
			if m.focusedPane == 0 {
				url := m.items[m.cursor]
				m.response = fmt.Sprintf("Checking %s...", url)
				return m, client.CheckServer(url)
			}

		case "n":
			if m.focusedPane == 0 {
				m.isTyping = true
				m.textInput.Focus()
				return m, textinput.Blink
			}
		
		case "tab", "left", "right":
			if m.focusedPane == 0 {
				m.focusedPane = 1
			} else {
				m.focusedPane = 0
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	header := titleStyle.Render("Welcome to lazy-curl!") + "\n"

	leftContent := ""
	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

        line := fmt.Sprintf("%s [%s] %s\n", cursor, checked, item)

		if m.cursor == i && m.focusedPane == 0 {
			leftContent += cursorStyle.Render(line) + "\n"
		} else {
			leftContent += itemStyle.Render(line) + "\n"
		}
	}

	if m.isTyping {
		leftContent += fmt.Sprintf("\nEnter new URL:\n%s\n(Press Enter to save, Esc to cancel)\n", m.textInput.View())
	} else {
		leftContent += "\nPress 'n' to add a new URL.\n"
	}

	rightContent := m.responseBody

	leftBox := baseBoxStyle.Copy().Width(40).Height(15)
	rightBox := baseBoxStyle.Copy().Width(50).Height(15)

	if m.focusedPane == 0 {
		leftBox = leftBox.BorderForeground(activeBorder)
		rightBox = rightBox.BorderForeground(inactiveBorder)
	} else {
		leftBox = leftBox.BorderForeground(inactiveBorder)
		rightBox = rightBox.BorderForeground(activeBorder)
	}

	renderLeft := leftBox.Render(leftContent)
	renderRight := rightBox.Render(rightContent)

	grid := lipgloss.JoinHorizontal(lipgloss.Top, renderLeft, renderRight)

	status := "\nStatus:" + statusStyle.Render(m.response) + "\n"

	helper := itemStyle.Render("\nPress j/k to move, Space to select, Enter to check, q to quit.\n")

	return header + grid + status + helper
}
