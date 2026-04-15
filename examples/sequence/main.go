package main

// A simple example illustrating how to run a series of commands in order.

import (
	"fmt"
	"os"
	"time"

	"github.com/carmel/go-tui"
)

type model struct{}

func (m model) Init() tui.Cmd {
	// A tui.Sequence is a command that runs a series of commands in
	// order. Contrast this with tui.Batch, which runs a series of commands
	// concurrently, with no order guarantees.
	return tui.Sequence(
		tui.Batch(
			tui.Sequence(
				SleepPrintln("1-1-1", 1000),
				SleepPrintln("1-1-2", 1000),
			),
			tui.Batch(
				SleepPrintln("1-2-1", 1500),
				SleepPrintln("1-2-2", 1250),
			),
		),
		tui.Println("2"),
		tui.Sequence(
			tui.Batch(
				SleepPrintln("3-1-1", 500),
				SleepPrintln("3-1-2", 1000),
			),
			tui.Sequence(
				SleepPrintln("3-2-1", 750),
				SleepPrintln("3-2-2", 500),
			),
		),
		tui.Quit,
	)
}

// print string after stopping for a certain period of time
func SleepPrintln(s string, milisecond int) tui.Cmd {
	printCmd := tui.Println(s)
	return func() tui.Msg {
		time.Sleep(time.Duration(milisecond) * time.Millisecond)
		return printCmd()
	}
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg.(type) {
	case tui.KeyPressMsg:
		return m, tui.Quit
	}
	return m, nil
}

func (m model) View() tui.View {
	return tui.NewView("")
}

func main() {
	if _, err := tui.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
