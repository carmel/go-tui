package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/textarea"
)

func main() {
	p := tui.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	textarea textarea.Model
}

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "Schnrr..."
	ti.ShowLineNumbers = true
	ti.DynamicHeight = true
	ti.MinHeight = 3
	ti.MaxHeight = 15
	ti.MaxContentHeight = 20
	ti.SetWidth(60)
	ti.SetVirtualCursor(false)
	ti.Focus()

	return model{textarea: ti}
}

func (m model) Init() tui.Cmd {
	return tui.Batch(textarea.Blink, tui.RequestBackgroundColor)
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	var cmds []tui.Cmd
	var cmd tui.Cmd

	switch msg := msg.(type) {
	case tui.BackgroundColorMsg:
		m.textarea.SetStyles(textarea.DefaultStyles(msg.IsDark()))
	case tui.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tui.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tui.Batch(cmds...)
}

func (m model) statusView() string {
	return fmt.Sprintf(
		"\nHeight: %d · Lines: %d · Cursor: (%d, %d) · Scroll: %.0f%%",
		m.textarea.Height(),
		m.textarea.LineCount(),
		m.textarea.Line(),
		m.textarea.Column(),
		m.textarea.ScrollPercent()*100,
	)
}

func (m model) View() tui.View {
	const gap = 1

	var c *tui.Cursor
	if !m.textarea.VirtualCursor() {
		c = m.textarea.Cursor()
		c.Y += gap
	}

	f := strings.Repeat("\n", gap)
	f += strings.Join([]string{
		m.textarea.View(),
		m.statusView(),
		"\n(ctrl+c to quit)",
	}, "\n")

	v := tui.NewView(f)
	v.Cursor = c
	return v
}
