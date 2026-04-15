package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"log"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/textinput"
)

func main() {
	p := tui.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
	quitting  bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return model{textInput: ti}
}

func (m model) Init() tui.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	var cmd tui.Cmd

	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		switch msg.String() {
		case "enter", "ctrl+c", "esc":
			m.quitting = true
			return m, tui.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() tui.View {
	var c *tui.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())
	if m.quitting {
		str += "\n"
	}

	v := tui.NewView(str)
	v.Cursor = c
	return v
}

func (m model) headerView() string { return "What’s your favorite Pokémon?\n" }
func (m model) footerView() string { return "\n(esc to quit)" }
