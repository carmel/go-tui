package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/carmel/go-tui"
)

const url = "https://charm.sh/"

type model struct {
	status int
	err    error
}

func checkServer() tui.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	}
	defer res.Body.Close() // nolint:errcheck

	return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tui.Cmd {
	return checkServer
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return m, tui.Quit

	case errMsg:
		m.err = msg
		return m, tui.Quit

	case tui.KeyPressMsg:
		if msg.Mod == tui.ModCtrl && msg.Code == 'c' {
			return m, tui.Quit
		}
	}

	return m, nil
}

func (m model) View() tui.View {
	if m.err != nil {
		return tui.NewView(fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err))
	}

	s := fmt.Sprintf("Checking %s ... ", url)
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	return tui.NewView("\n" + s + "\n\n")
}

func main() {
	if _, err := tui.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
