package main

import (
	"cmp"
	"fmt"
	"os"
	"os/exec"

	"github.com/carmel/go-tui"
)

type editorFinishedMsg struct{ err error }

func openEditor() tui.Cmd {
	editor := cmp.Or(os.Getenv("EDITOR"), "vim")
	c := exec.Command(editor) //nolint:gosec
	return tui.ExecProcess(c, func(err error) tui.Msg {
		return editorFinishedMsg{err}
	})
}

type model struct {
	altscreenActive bool
	err             error
}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		switch msg.String() {
		case "a":
			m.altscreenActive = !m.altscreenActive
			return m, nil
		case "e":
			return m, openEditor()
		case "ctrl+c", "q":
			return m, tui.Quit
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tui.Quit
		}
	}
	return m, nil
}

func (m model) View() tui.View {
	if m.err != nil {
		v := tui.NewView("Error: " + m.err.Error() + "\n")
		v.AltScreen = m.altscreenActive
		return v
	}
	v := tui.NewView("Press 'e' to open your EDITOR.\nPress 'a' to toggle the altscreen\nPress 'q' to quit.\n")
	v.AltScreen = m.altscreenActive
	return v
}

func main() {
	m := model{}
	if _, err := tui.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
