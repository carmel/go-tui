package main

import (
	"log"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
)

var body = lipgloss.NewStyle().Padding(1, 2)

type model struct {
	value int
	width int
	state tui.ProgressBarState
}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.WindowSizeMsg:
		m.width = msg.Width
	case tui.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tui.Quit
		case "up", "k":
			if m.value < 100 {
				m.value += 10
			}
		case "down", "j":
			if m.value > 0 {
				m.value -= 10
			}
		case "left", "h":
			if m.state > 0 {
				m.state--
			}
		case "right", "l":
			if m.state < 4 {
				m.state++
			}
		}
	}
	return m, nil
}

func (m model) View() tui.View {
	s := body.Width(m.width - body.GetHorizontalPadding()).Render(
		"This demo requires a terminal emulator that supports an indeterminate progress bar, such a Windows Terminal or Ghostty. In other terminals (including tmux in a supporting terminal) nothing will happen.\n\nPress up/down to change value, left/right to change state, q to quit.",
	)
	v := tui.NewView(s)
	v.ProgressBar = tui.NewProgressBar(m.state, m.value)
	return v
}

func main() {
	p := tui.NewProgram(model{value: 50, state: tui.ProgressBarIndeterminate})
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
