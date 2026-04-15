package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/carmel/go-tui"
)

type model struct {
	quitting   bool
	suspending bool
}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.ResumeMsg:
		m.suspending = false
		return m, nil
	case tui.KeyPressMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tui.Quit
		case "ctrl+c":
			m.quitting = true
			return m, tui.Interrupt
		case "ctrl+z":
			m.suspending = true
			return m, tui.Suspend
		}
	}
	return m, nil
}

func (m model) View() tui.View {
	if m.suspending || m.quitting {
		return tui.NewView("")
	}

	return tui.NewView("\nPress ctrl-z to suspend, ctrl+c to interrupt, q, or esc to exit\n")
}

func main() {
	if _, err := tui.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		if errors.Is(err, tui.ErrInterrupted) {
			os.Exit(130)
		}
		os.Exit(1)
	}
}
