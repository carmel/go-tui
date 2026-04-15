package main

// A simple program that queries and displays the window-size.

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
		return m, tui.RequestWindowSize

	case tui.WindowSizeMsg:
		return m, tui.Printf("The window size is: %dx%d", msg.Width, msg.Height)
	}

	return m, nil
}

func (m model) View() tui.View {
	return tui.NewView("\nWhen you're done press q to quit.\nPress any other key to query the window-size.\n")
}
