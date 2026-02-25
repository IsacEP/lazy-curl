package main

import (
	"fmt"
	"os"

	"github.com/IsacEP/lazy-curl/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)



func main() {
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There was an Error")
		os.Exit(1)
	}
}
