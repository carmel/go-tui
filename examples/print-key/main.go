package main

import (
	"log"

	"github.com/carmel/go-tui"
)

type model struct{}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyboardEnhancementsMsg:
		return m, tui.Printf("Keyboard enhancements: EventTypes: %v\n",
			msg.SupportsEventTypes())
	case tui.KeyMsg:
		key := msg.Key()
		switch msg := msg.(type) {
		case tui.KeyPressMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tui.Quit
			}
		}
		format := "(%T) You pressed: %s"
		args := []any{msg, msg.String()}
		if len(key.Text) > 0 {
			format += " (text: %q)"
			args = append(args, key.Text)
		}
		return m, tui.Printf(format, args...)
	}
	return m, nil
}

func (m model) View() tui.View {
	v := tui.NewView("Press any key to see its details printed to the terminal. Press 'ctrl+c' to quit.")
	v.KeyboardEnhancements.ReportEventTypes = true
	return v
}

func main() {
	p := tui.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
	}
}
