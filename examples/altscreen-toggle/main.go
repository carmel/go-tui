package main

import (
	"fmt"
	"os"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
)

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	altscreen  bool
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
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tui.Quit
		case "ctrl+z":
			m.suspending = true
			return m, tui.Suspend
		case "space":
			var cmd tui.Cmd
			m.altscreen = !m.altscreen
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() tui.View {
	if m.suspending {
		v := tui.NewView("")
		v.AltScreen = m.altscreen
		return v
	}

	if m.quitting {
		v := tui.NewView("Bye!\n")
		v.AltScreen = m.altscreen
		return v
	}

	const (
		altscreenMode = " altscreen mode "
		inlineMode    = " inline mode "
	)

	var mode string
	if m.altscreen {
		mode = altscreenMode
	} else {
		mode = inlineMode
	}

	v := tui.NewView(fmt.Sprintf("\n\n  You're in %s\n\n\n", keywordStyle.Render(mode)) +
		helpStyle.Render("  space: switch modes • ctrl-z: suspend • q: exit\n"))
	v.AltScreen = m.altscreen
	return v
}

func main() {
	if _, err := tui.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
