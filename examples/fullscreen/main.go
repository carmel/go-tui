package main

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"fmt"
	"log"
	"time"

	"github.com/carmel/go-tui"
)

type model int

type tickMsg time.Time

func main() {
	p := tui.NewProgram(model(5))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tui.Cmd {
	return tick()
}

func (m model) Update(message tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := message.(type) {
	case tui.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tui.Quit
		}

	case tickMsg:
		m--
		if m <= 0 {
			return m, tui.Quit
		}
		return m, tick()
	}

	return m, nil
}

func (m model) View() tui.View {
	v := tui.NewView(fmt.Sprintf("\n\n     Hi. This program will exit in %d seconds...", m))
	v.AltScreen = true
	return v
}

func tick() tui.Cmd {
	return tui.Tick(time.Second, func(t time.Time) tui.Msg {
		return tickMsg(t)
	})
}
