package main

import (
	"fmt"
	"os"

	"github.com/carmel/go-tui"
)

type model struct {
	cursor tui.Cursor
	blink  bool
}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tui.Quit
		case "h", "left":
			m.cursor.Shape--
			if m.cursor.Shape < tui.CursorBlock {
				m.cursor.Shape = tui.CursorBar
			}
		case "l", "right":
			m.cursor.Shape++
			if m.cursor.Shape > tui.CursorBar {
				m.cursor.Shape = tui.CursorBlock
			}
		}
	}
	m.blink = !m.blink
	return m, nil
}

func (m model) View() tui.View {
	v := tui.NewView("Press left/right to change the cursor style, q or ctrl+c to quit." +
		"\n\n" +
		"  <- This is the cursor (a " + m.describeCursor() + ")")
	c := tui.NewCursor(0, 2)
	c.Shape = m.cursor.Shape
	c.Blink = m.blink
	v.Cursor = c
	return v
}

func (m model) describeCursor() string {
	var adj, noun string

	if m.blink {
		adj = "blinking"
	} else {
		adj = "steady"
	}

	switch m.cursor.Shape {
	case tui.CursorBlock:
		noun = "block"
	case tui.CursorUnderline:
		noun = "underline"
	case tui.CursorBar:
		noun = "bar"
	}

	return fmt.Sprintf("%s %s", adj, noun)
}

func main() {
	p := tui.NewProgram(model{blink: true})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
