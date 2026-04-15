package main

// A simple program that handled losing and acquiring focus.

import (
	"log"

	"github.com/carmel/go-tui"
)

func main() {
	p := tui.NewProgram(model{
		focused:   true,
		reporting: true,
	})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	focused   bool
	reporting bool
}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.FocusMsg:
		m.focused = true
	case tui.BlurMsg:
		m.focused = false
	case tui.KeyPressMsg:
		switch msg.String() {
		case "t":
			m.reporting = !m.reporting
		case "ctrl+c", "q":
			return m, tui.Quit
		}
	}

	return m, nil
}

func (m model) View() tui.View {
	s := "Hi. Focus report is currently "
	if m.reporting {
		s += "enabled"
	} else {
		s += "disabled"
	}
	s += ".\n\n"

	if m.reporting {
		if m.focused {
			s += "This program is currently focused!"
		} else {
			s += "This program is currently blurred!"
		}
	}
	v := tui.NewView(s + "\n\nTo quit sooner press ctrl-c, or t to toggle focus reporting...\n")
	v.ReportFocus = m.reporting
	return v
}
