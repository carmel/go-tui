package main

// An example illustrating how to pipe in data to a Bubble Tea application.
// More so, this serves as proof that Bubble Tea will automatically listen for
// keystrokes when input is not a TTY, such as when data is piped or redirected
// in.

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/textinput"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		fmt.Println("Try piping in some text.")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			fmt.Println("Error getting input:", err)
			os.Exit(1)
		}
	}

	model := newModel(strings.TrimSpace(b.String()))

	if _, err := tui.NewProgram(model).Run(); err != nil {
		fmt.Println("Couldn't start program:", err)
		os.Exit(1)
	}
}

type model struct {
	userInput textinput.Model
}

func newModel(initialValue string) (m model) {
	i := textinput.New()
	i.Prompt = ""

	s := i.Styles()
	s.Cursor.Color = lipgloss.Color("63")
	i.SetStyles(s)

	i.SetWidth(48)
	i.SetValue(initialValue)
	i.CursorEnd()
	i.Focus()

	m.userInput = i
	return m
}

func (m model) Init() tui.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	if key, ok := msg.(tui.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc", "enter":
			return m, tui.Quit
		}
	}

	var cmd tui.Cmd
	m.userInput, cmd = m.userInput.Update(msg)
	return m, cmd
}

func (m model) View() tui.View {
	return tui.NewView(fmt.Sprintf(
		"\nYou piped in: %s\n\nPress ^C to exit",
		m.userInput.View(),
	))
}
