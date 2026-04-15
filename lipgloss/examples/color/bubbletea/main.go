package main

import (
	"fmt"
	"os"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
)

// Style definitions.
type styles struct {
	frame,
	paragraph,
	text,
	keyword,
	activeButton,
	inactiveButton lipgloss.Style
}

// Styles are initialized based on the background color of the terminal.
func newStyles(backgroundIsDark bool) (s *styles) {
	s = new(styles)

	// Create a new helper function for choosing either a light or dark color
	// based on the detected background color.
	lightDark := lipgloss.LightDark(backgroundIsDark)

	// Define some styles. adaptive.Color() can be used to choose the
	// appropriate light or dark color based on the detected background color.
	s.frame = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightDark(
			lipgloss.Color("#C5ADF9"),
			lipgloss.Color("#864EFF"))).
		Padding(1, 3).
		Margin(1, 3)
	s.paragraph = lipgloss.NewStyle().
		Width(40).
		MarginBottom(1).
		Align(lipgloss.Center)
	s.text = lipgloss.NewStyle().
		Foreground(lightDark(
			lipgloss.Color("#696969"),
			lipgloss.Color("#bdbdbd")))
	s.keyword = lipgloss.NewStyle().
		Foreground(lightDark(
			lipgloss.Color("#37CD96"),
			lipgloss.Color("#22C78A"))).
		Bold(true)

	s.activeButton = lipgloss.NewStyle().
		Padding(0, 3).
		Background(lipgloss.Color("#FF6AD2")).
		Foreground(lipgloss.Color("#FFFCC2"))
	s.inactiveButton = s.activeButton.
		Background(lightDark(
			lipgloss.Color("#988F95"),
			lipgloss.Color("#978692"))).
		Foreground(lightDark(
			lipgloss.Color("#FDFCE3"),
			lipgloss.Color("#FBFAE7")))
	return s
}

type model struct {
	styles  *styles
	yes     bool
	chosen  bool
	aborted bool
}

func (m model) Init() tui.Cmd {
	// Query for the background color on start.
	m.yes = true
	return tui.RequestBackgroundColor
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {

	// Bubble Tea automatically detects the background color on start. We
	// listen for the response here, then initialize our styles accordingly.
	case tui.BackgroundColorMsg:
		m.styles = newStyles(msg.IsDark())
		return m, nil

	case tui.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.aborted = true
			return m, tui.Quit
		case "enter":
			m.chosen = true
			return m, tui.Quit
		case "left", "right", "h", "l":
			m.yes = !m.yes
		case "y":
			m.yes = true
			m.chosen = true
			return m, tui.Quit
		case "n":
			m.yes = false
			m.chosen = true
			return m, tui.Quit
		}
	}

	return m, nil
}

func (m model) View() tui.View {
	var v tui.View
	if m.styles == nil {
		// We haven't received tui.BackgroundColorMsg yet. Don't worry, it'll
		// be here in a flash.
		return v
	}
	if m.chosen || m.aborted {
		// We're about to exit, so wipe the UI.
		return v
	}

	var (
		s = m.styles
		y = "Yes"
		n = "No"
	)

	if m.yes {
		y = s.activeButton.Render(y)
		n = s.inactiveButton.Render(n)
	} else {
		y = s.inactiveButton.Render(y)
		n = s.activeButton.Render(n)
	}

	content := s.frame.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			s.paragraph.Render(
				s.text.Render("Are you sure you want to eat that ")+
					s.keyword.Render("moderatly ripe")+
					s.text.Render(" banana?"),
			),
			y+"  "+n,
		),
	)
	v.SetContent(content)
	return v
}

func main() {
	m, err := tui.NewProgram(model{}).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uh oh: %v", err)
		os.Exit(1)
	}

	if m := m.(model); m.chosen {
		if m.yes {
			fmt.Println("Are you sure? It's not ripe yet.")
		} else {
			fmt.Println("Well, alright. It was probably good, though.")
		}
	}
}
