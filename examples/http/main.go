package main

// A simple program that makes a GET request and prints the response status.

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/carmel/go-tui"
)

const url = "https://charm.sh/"

type model struct {
	status int
	err    error
}

type statusMsg int

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func main() {
	p := tui.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tui.Cmd {
	return checkServer
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tui.Quit
		default:
			return m, nil
		}

	case statusMsg:
		m.status = int(msg)
		return m, tui.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() tui.View {
	s := fmt.Sprintf("Checking %s...", url)
	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	} else if m.status != 0 {
		s += fmt.Sprintf("%d %s", m.status, http.StatusText(m.status))
	}
	return tui.NewView(s + "\n")
}

func checkServer() tui.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	}
	defer res.Body.Close() // nolint:errcheck

	return statusMsg(res.StatusCode)
}
