package main

import (
	"strings"
	"time"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/progress"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type progressMsg float64

type progressErrMsg struct{ err error }

func finalPause() tui.Cmd {
	return tui.Tick(time.Millisecond*750, func(_ time.Time) tui.Msg {
		return nil
	})
}

type model struct {
	pw       *progressWriter
	progress progress.Model
	err      error
}

func (m model) Init() tui.Cmd {
	return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyPressMsg:
		return m, tui.Quit

	case tui.WindowSizeMsg:
		m.progress.SetWidth(msg.Width - padding*2 - 4)
		if m.progress.Width() > maxWidth {
			m.progress.SetWidth(maxWidth)
		}
		return m, nil

	case progressErrMsg:
		m.err = msg.err
		return m, tui.Quit

	case progressMsg:
		var cmds []tui.Cmd

		if msg >= 1.0 {
			cmds = append(cmds, tui.Sequence(finalPause(), tui.Quit))
		}

		cmds = append(cmds, m.progress.SetPercent(float64(msg)))
		return m, tui.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var cmd tui.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() tui.View {
	if m.err != nil {
		return tui.NewView("Error downloading: " + m.err.Error() + "\n")
	}

	pad := strings.Repeat(" ", padding)
	return tui.NewView("\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit"))
}
