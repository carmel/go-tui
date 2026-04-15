package main

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"log"
	"strings"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/textarea"
)

func main() {
	p := tui.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error

type model struct {
	textarea textarea.Model
	err      error
}

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.SetVirtualCursor(false)
	ti.SetStyles(textarea.DefaultStyles(true)) // default to dark styles.
	ti.Focus()

	return model{
		textarea: ti,
		err:      nil,
	}
}

func (m model) Init() tui.Cmd {
	return tui.Batch(textarea.Blink, tui.RequestBackgroundColor)
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	var cmds []tui.Cmd
	var cmd tui.Cmd

	switch msg := msg.(type) {
	case tui.BackgroundColorMsg:
		// Update styling now that we know the background color.
		m.textarea.SetStyles(textarea.DefaultStyles(msg.IsDark()))

	case tui.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case "ctrl+c":
			return m, tui.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

		// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tui.Batch(cmds...)
}

func (m model) headerView() string {
	return "Tell me a story.\n"
}

func (m model) View() tui.View {
	const (
		footer = "\n(ctrl+c to quit)\n"
	)

	var c *tui.Cursor
	if !m.textarea.VirtualCursor() {
		c = m.textarea.Cursor()

		if c != nil {
			// Set the y offset of the cursor based on the position of the textarea
			// in the application.
			offset := lipgloss.Height(m.headerView())
			c.Y += offset
		}
	}

	f := strings.Join([]string{
		m.headerView(),
		m.textarea.View(),
		footer,
	}, "\n")

	v := tui.NewView(f)
	v.Cursor = c
	return v
}
