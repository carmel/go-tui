package main

import (
	"fmt"
	"os"

	"github.com/carmel/go-tui"
)

type model bool

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	if _, ok := msg.(tui.KeyPressMsg); ok {
		m = true
		return m, tui.Quit
	}
	return m, nil
}

func (m model) View() tui.View {
	if m {
		return tui.NewView("")
	}
	return tui.NewView("Press any key to quit.\n(When this program quits, it will vanish without a trace.)")
}

func main() {
	p := tui.NewProgram(model(false))
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Oh no:", err)
	}
}
