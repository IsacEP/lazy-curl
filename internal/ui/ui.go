package ui

import (
	"fmt"
	"strings"

	"github.com/IsacEP/lazy-curl/internal/client"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).MarginBottom(1)

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)

	endpointstyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0"))

	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7CCB")).Italic(true)

	activeBorder = lipgloss.Color("#ff5ec9")
	inactiveBorder = lipgloss.Color("#91579b")

	baseBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

type Model struct {
	cursor int
	endpoints []string
	selected map[int]struct{}
	response string
	textInput textinput.Model
	isTyping  bool
	focusedPane int
	responseBody string
	baseURL string
	baseURLinput textinput.Model
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter New Endpoint: "
	ti.CharLimit = 156
	ti.Width = 50

	urlti := textinput.New()
	urlti.Placeholder = "Please Enter Base URL: "
	urlti.CharLimit = 156
	urlti.Width = 50

	return Model {
		endpoints:     []string{"TEST"},
		selected:  make(map[int]struct{}),
		response:  "Ready to use server",
		textInput: ti,
		isTyping:  false,
		focusedPane: 0,
		responseBody: "response body",
		baseURL: "",
		baseURLinput: urlti,
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
				if m.focusedPane == 1 {
					val := m.textInput.Value()
					if val != "" {
						m.endpoints = append(m.endpoints, val)
						m.textInput.Reset()
					}
					m.isTyping = false
					m.textInput.Blur()
					return m, nil
				}
				if m.focusedPane == 0 {
					val := m.baseURLinput.Value()
					if val != "" {
						m.baseURL = val
					}
					m.isTyping = false
					m.baseURLinput.Blur()
					return m, nil
				}
			case "esc":
				m.isTyping = false
				m.textInput.Blur()
				return m, nil
			}

			if m.focusedPane == 1 {
				m.textInput, cmd = m.textInput.Update(msg)
			} else if m.focusedPane == 0 {
				m.baseURLinput, cmd = m.baseURLinput.Update(msg)
			}
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.focusedPane == 1 && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.focusedPane == 1 && m.cursor < len(m.endpoints)-1 {
				m.cursor++
			}

		case " ":
			if m.focusedPane == 1 {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}

		case "enter":
			if m.focusedPane == 1 {
				rawEndpoint := m.endpoints[m.cursor]
				path := rawEndpoint
				parts := strings.Fields(rawEndpoint)
				if (len(parts)) > 1 {
					path = parts[1]
				}
				cleanBase := strings.TrimSuffix(m.baseURL, "/")
				cleanPath := strings.TrimPrefix(path, "/")

				fullURL := cleanBase + "/" + cleanPath
				if !strings.HasPrefix(fullURL, "http://") && !strings.HasPrefix(fullURL, "https://") {
					fullURL = "http://" + fullURL
				}
				m.response = fmt.Sprintf("Checking %s...", fullURL)
				return m, client.CheckServer(fullURL)
			}

		case "n":
			if m.focusedPane == 1 {
				m.isTyping = true
				m.textInput.Focus()
				return m, textinput.Blink
			}
		case "u":
			if m.focusedPane == 0 {
				m.isTyping = true
				m.baseURLinput.Focus()
				return m, textinput.Blink
			}
		case "tab":
			m.focusedPane = (m.focusedPane + 1) % 3
		case "left":
			m.focusedPane--
			if m.focusedPane < 0 {
				m.focusedPane = 2
			}
		case "right":
			m.focusedPane = (m.focusedPane + 1) % 3
		}
	}
	return m, nil
}

func (m Model) View() string {
	header := titleStyle.Render("Welcome to lazy-curl!") + "\n"
	leftContent := ""
	for i, item := range m.endpoints {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

        line := fmt.Sprintf("%s [%s] %s\n", cursor, checked, item)

		if m.cursor == i && m.focusedPane == 1 {
			leftContent += cursorStyle.Render(line) + "\n"
		} else {
			leftContent += endpointstyle.Render(line) + "\n"
		}
	}

	if m.isTyping {
		leftContent += fmt.Sprintf("\nEnter new URL:\n%s\n(Press Enter to save, Esc to cancel)\n", m.textInput.View())
	} else {
		leftContent += "\nPress 'n' to add a new Endpoint.\n"
	}

	rightContent := m.responseBody

	topContent := ""
	if m.isTyping {
		topContent += fmt.Sprintf("Set Base URL:\n%s\n", m.baseURLinput.View())
	} else {
		if m.baseURL != "" {
			topContent += m.baseURL
		} else {
			topContent += "\nPress 'u' to add a new base URL"
		}
	}
	topBox := baseBoxStyle.Copy().Width(94).Height(5)
	leftBox := baseBoxStyle.Copy().Width(40).Height(15)
	rightBox := baseBoxStyle.Copy().Width(50).Height(15)

	if m.focusedPane == 0 {
		topBox = topBox.BorderForeground(activeBorder)
	} else {
		topBox = topBox.BorderForeground(inactiveBorder)
	}
	if m.focusedPane == 1 {
		leftBox = leftBox.BorderForeground(activeBorder)
	} else {
		leftBox = leftBox.BorderForeground(inactiveBorder)
	}

	if m.focusedPane == 2 {
		rightBox = rightBox.BorderForeground(activeBorder)
	} else {
		rightBox = rightBox.BorderForeground(inactiveBorder)
	}

	renderTop := topBox.Render(topContent)
	renderLeft := leftBox.Render(leftContent)
	renderRight := rightBox.Render(rightContent)

	bottomHalf := lipgloss.JoinHorizontal(lipgloss.Top, renderLeft, renderRight)

	grid := lipgloss.JoinVertical(lipgloss.Left, renderTop, bottomHalf)

	status := "\nStatus:" + statusStyle.Render(m.response) + "\n"

	helper := endpointstyle.Render("\nPress j/k to move, Space to select, Enter to check, q to quit.\n")

	return header + grid + status + helper
}
