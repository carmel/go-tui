package main

// A program demonstrating how to use the WithFilter option to intercept events.

import (
	"fmt"
	"log"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/help"
	"github.com/carmel/go-tui/key"
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/textarea"
)

var (
	choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))
	saveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	quitViewStyle = lipgloss.NewStyle().Padding(1, 3).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
)

func main() {
	p := tui.NewProgram(initialModel(), tui.WithFilter(filter))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func filter(teaModel tui.Model, msg tui.Msg) tui.Msg {
	if _, ok := msg.(tui.QuitMsg); !ok {
		return msg
	}

	m := teaModel.(model)
	if m.hasChanges {
		return nil
	}

	return msg
}

type model struct {
	textarea   textarea.Model
	help       help.Model
	keymap     keymap
	saveText   string
	hasChanges bool
	quitting   bool
}

type keymap struct {
	save key.Binding
	quit key.Binding
}

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "Only the best words"
	ti.Focus()

	return model{
		textarea: ti,
		help:     help.New(),
		keymap: keymap{
			save: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "save"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
}

func (m model) Init() tui.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	if m.quitting {
		return m.updatePromptView(msg)
	}

	return m.updateTextView(msg)
}

func (m model) updateTextView(msg tui.Msg) (tui.Model, tui.Cmd) {
	var cmds []tui.Cmd
	var cmd tui.Cmd

	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		m.saveText = ""
		switch {
		case key.Matches(msg, m.keymap.save):
			m.saveText = "Changes saved!"
			m.hasChanges = false
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tui.Quit
		case len(msg.Text) > 0:
			m.saveText = ""
			m.hasChanges = true
			fallthrough
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tui.Batch(cmds...)
}

func (m model) updatePromptView(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		// For simplicity's sake, we'll treat any key besides "y" as "no"
		if key.Matches(msg, m.keymap.quit) || msg.String() == "y" {
			m.hasChanges = false
			return m, tui.Quit
		}
		m.quitting = false
	}

	return m, nil
}

func (m model) View() tui.View {
	if m.quitting {
		if m.hasChanges {
			text := lipgloss.JoinHorizontal(lipgloss.Top, "You have unsaved changes. Quit without saving?", choiceStyle.Render("[yN]"))
			return tui.NewView(quitViewStyle.Render(text))
		}
		return tui.NewView("Very important. Thank you.\n")
	}

	helpView := m.help.ShortHelpView([]key.Binding{
		m.keymap.save,
		m.keymap.quit,
	})

	return tui.NewView(fmt.Sprintf(
		"Type some important things.\n%s\n %s\n %s",
		m.textarea.View(),
		saveTextStyle.Render(m.saveText),
		helpView,
	) + "\n\n")
}
