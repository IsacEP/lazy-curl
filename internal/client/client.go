package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerStatusMsg struct {
	URL string
	Status int
	Err error
	Body string
}

func CheckServer(url string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 5 * time.Second}
		resp, err := c.Get(url)
		
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, bodyBytes, "", "  ")

		bodyString := ""
		if err == nil {
			bodyString = prettyJSON.String()
		} else {
			bodyString = string(bodyBytes) 
		}

		return ServerStatusMsg{URL: url, Status: resp.StatusCode, Body: bodyString}
	}
}

func SendRequest(url string, method string, data string) tea.Cmd {
	switch method {
	case "GET":
		return get(url)
	case "POST":
		return post(url, data)
	case "PUT":
		return put(url, data)
	case "PATCH":
		return patch(url, data)
	case "DELETE":
		delete(url, data)
	}
}

func delete(url, data string) tea.Cmd {
	panic("unimplemented")
}

func patch(url, data string) tea.Cmd {
	panic("unimplemented")
}

func put(url, data string) tea.Cmd {
	panic("unimplemented")
}

func post(url, data string) tea.Cmd {
	panic("unimplemented")
}

func get(url string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 5 * time.Second}
		resp, err := c.Get(url)
		
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, bodyBytes, "", "  ")

		bodyString := ""
		if err == nil {
			bodyString = prettyJSON.String()
		} else {
			bodyString = string(bodyBytes) 
		}

		return ServerStatusMsg{URL: url, Status: resp.StatusCode, Body: bodyString}
	}
}