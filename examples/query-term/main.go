// This example uses a textinput to send the terminal ANSI sequences to query
// it for capabilities.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/textinput"
)

func newModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)
	ti.SetVirtualCursor(false)
	return model{input: ti}
}

type model struct {
	input textinput.Model
	err   error
}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	var cmds []tui.Cmd

	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		m.err = nil
		switch msg.String() {
		case "ctrl+c":
			return m, tui.Quit
		case "enter":
			// Write the sequence to the terminal.
			val := m.input.Value()
			val = "\"" + val + "\""

			// Unescape the sequence.
			seq, err := strconv.Unquote(val)
			if err != nil {
				m.err = err
				return m, nil
			}

			if !strings.HasPrefix(seq, "\x1b") {
				m.err = fmt.Errorf("sequence is not an ANSI escape sequence")
				return m, nil
			}

			m.input.SetValue("")

			// Write the sequence to the terminal.
			return m, func() tui.Msg {
				io.WriteString(os.Stdout, seq)
				return nil
			}
		}
	default:
		_, typ, ok := strings.Cut(fmt.Sprintf("%T", msg), ".")
		if ok && unicode.IsUpper(rune(typ[0])) {
			// Only log messages that are exported types.
			cmds = append(cmds, tui.Printf("Received message: %T %+v", msg, msg))
		}
	}

	var cmd tui.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tui.Batch(cmds...)
}

func (m model) View() tui.View {
	var s strings.Builder
	s.WriteString(m.input.View())
	if m.err != nil {
		s.WriteString("\n\nError: " + m.err.Error())
	}
	s.WriteString("\n\nPress ctrl+c to quit, enter to write the sequence to terminal")
	v := tui.NewView(s.String())
	v.Cursor = m.input.Cursor()
	return v
}

func main() {
	p := tui.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
