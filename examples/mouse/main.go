package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"log"

	"github.com/carmel/go-tui"
)

func main() {
	p := tui.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct{}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tui.Quit
		}

	case tui.MouseMsg:
		mouse := msg.Mouse()
		return m, tui.Printf("(X: %d, Y: %d) %s", mouse.X, mouse.Y, mouse)
	}

	return m, nil
}

func (m model) View() tui.View {
	v := tui.NewView("Do mouse stuff. When you're done press q to quit.\n")
	v.MouseMode = tui.MouseModeAllMotion
	return v
}
