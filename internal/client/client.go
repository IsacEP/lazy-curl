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
		return doGet(url)
	case "POST":
		return doPost(url, data)
	case "PUT":
		return doPut(url, data)
	case "PATCH":
		return doPatch(url, data)
	case "DELETE":
		return doDelete(url, data)
	default:
		return doGet(url)
	}
}

func doGet(url string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		return executeRequest(req, url)
	}
}

func doPost(url string, data string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		req.Header.Set("Content-Type", "application/json")
		return executeRequest(req, url)
	}
}

func doPut(url string, data string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		req.Header.Set("Content-Type", "application/json")
		return executeRequest(req, url)
	}
}

func doPatch(url string, data string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		req.Header.Set("Content-Type", "application/json")
		return executeRequest(req, url)
	}
}

func doDelete(url string, data string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return ServerStatusMsg{URL: url, Err: err}
		}
		req.Header.Set("Content-Type", "application/json")
		return executeRequest(req, url)
	}
}

func executeRequest(req *http.Request, url string) tea.Msg {
	c := &http.Client{Timeout: 5 * time.Second}
	resp, err := c.Do(req)

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