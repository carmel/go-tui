package main

import (
	"fmt"
	"os"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/textinput"
)

type model struct {
	input textinput.Model
	width int
}

var _ tui.Model = model{}

// Init implements tui.Model.
func (m model) Init() tui.Cmd {
	return m.input.Focus()
}

// Update implements tui.Model.
func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	var cmd tui.Cmd
	switch msg := msg.(type) {
	case tui.WindowSizeMsg:
		m.width = msg.Width
	case tui.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tui.Quit
		case "enter":
			input := m.input.Value()
			m.input.Reset()
			return m, tui.RequestCapability(input)
		}
	case tui.CapabilityMsg:
		return m, tui.Printf("Got capability: %s", msg)
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View implements tui.Model.
func (m model) View() tui.View {
	w := min(m.width, 60)

	instructions := lipgloss.NewStyle().
		Width(w).
		Render("Query for terminal capabilities. You can enter things like 'TN', 'RGB', 'cols', and so on. This will not work in all terminals and multiplexers.")

	return tui.NewView("\n" + instructions + "\n\n" +
		m.input.View() +
		"\n\nPress enter to request capability, or ctrl+c to quit.")
}

func main() {
	m := model{}
	m.input = textinput.New()
	m.input.Placeholder = "Enter capability name to request"
	m.input.Focus()

	if _, err := tui.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Uh oh:", err)
		os.Exit(1)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
