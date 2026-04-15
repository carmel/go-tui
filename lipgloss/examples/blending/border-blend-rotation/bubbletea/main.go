package main

import (
	"fmt"
	"os"
	"time"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
)

const (
	borderRotationFPS   = 15
	borderRotationSteps = 5
)

type borderRotationTickMsg struct {
	Value int
}

func borderRotationTick(current int) tui.Cmd {
	return tui.Tick(time.Second/time.Duration(borderRotationFPS), func(_ time.Time) tui.Msg {
		return borderRotationTickMsg{Value: current + borderRotationSteps}
	})
}

type model struct {
	borderRotation int
}

func (m model) Init() tui.Cmd {
	return borderRotationTick(0)
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tui.Quit
		}
	case borderRotationTickMsg:
		m.borderRotation = msg.Value
		return m, borderRotationTick(msg.Value)
	}

	return m, nil
}

func (m model) View() tui.View {
	v := tui.NewView(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForegroundBlend(
			lipgloss.Color("#00FA68"),
			lipgloss.Color("#9900FF"),
			lipgloss.Color("#ED5353"),
			lipgloss.Color("#9900FF"),
			lipgloss.Color("#00FA68"),
		).
		BorderForegroundBlendOffset(m.borderRotation).
		Width(60).
		Height(15).
		Render("Hello, world!"))
	v.AltScreen = true
	return v
}

func main() {
	_, err := tui.NewProgram(model{}).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uh oh: %v", err)
		os.Exit(1)
	}
}
