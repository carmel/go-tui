package main

// A simple example illustrating how to set a window title.

import (
	"fmt"
	"os"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
)

const windowTitle = "Hello, Bubble Tea"

type model struct{}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg.(type) {
	case tui.KeyPressMsg:
		return m, tui.Quit
	}
	return m, nil
}

func (m model) View() tui.View {
	wrap := lipgloss.NewStyle().Width(78).Render
	v := tui.NewView(wrap("The window title has been set to '"+windowTitle+"'. It will be cleared on exit.") +
		"\n\nPress any key to quit.")
	v.WindowTitle = windowTitle
	return v
}

func main() {
	if _, err := tui.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
