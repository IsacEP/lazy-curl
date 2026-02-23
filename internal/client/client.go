package client

import (
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerStatusMsg struct {
	URL string
	Status int
	Err error
}

func CheckServer(url string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 5 * time.Second}
		resp, err := c.Get(url)

		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		defer resp.Body.Close()

		return ServerStatusMsg{URL: url, Status: resp.StatusCode}
	}
}