package main

// A simple program demonstrating the paginator component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/paginator"

	"github.com/carmel/go-tui"
)

type styles struct {
	activeDot   lipgloss.Style
	inactiveDot lipgloss.Style
}

func newStyles(bgIsDark bool) (s styles) {
	lightDark := lipgloss.LightDark(bgIsDark)

	s.activeDot = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("235"), lipgloss.Color("252"))).SetString("•")
	s.inactiveDot = s.activeDot.Foreground(lightDark(lipgloss.Color("250"), lipgloss.Color("238"))).SetString("•")
	return s
}

type model struct {
	items     []string
	paginator paginator.Model
}

func newModel() model {
	var items []string
	for i := 1; i < 101; i++ {
		text := fmt.Sprintf("Item %d", i)
		items = append(items, text)
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.SetTotalPages(len(items))

	m := model{
		paginator: p,
		items:     items,
	}

	m.updateStyles(true) // default to dark styles
	return m
}

func (m *model) updateStyles(isDark bool) {
	styles := newStyles(isDark)
	m.paginator.ActiveDot = styles.activeDot.String()
	m.paginator.InactiveDot = styles.inactiveDot.String()
}

func (m model) Init() tui.Cmd {
	return tui.RequestBackgroundColor
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	var cmd tui.Cmd
	switch msg := msg.(type) {
	case tui.BackgroundColorMsg:
		m.updateStyles(msg.IsDark())
		return m, nil
	case tui.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tui.Quit
		}
	}
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m model) View() tui.View {
	var b strings.Builder
	b.WriteString("\n  Paginator Example\n\n")
	start, end := m.paginator.GetSliceBounds(len(m.items))
	for _, item := range m.items[start:end] {
		b.WriteString("  • " + item + "\n\n")
	}
	b.WriteString("  " + m.paginator.View())
	b.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	return tui.NewView(b.String())
}

func main() {
	p := tui.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
