package main

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/carmel/go-tui"
	"github.com/carmel/go-tui/spinner"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type responseMsg struct{}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func listenForActivity(sub chan struct{}) tui.Cmd {
	return func() tui.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100)) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tui.Cmd {
	return func() tui.Msg {
		return responseMsg(<-sub)
	}
}

type model struct {
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	spinner   spinner.Model
	quitting  bool
}

func (m model) Init() tui.Cmd {
	return tui.Batch(
		m.spinner.Tick,
		listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),   // wait for activity
	)
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg.(type) {
	case tui.KeyPressMsg:
		m.quitting = true
		return m, tui.Quit
	case responseMsg:
		m.responses++                    // record external activity
		return m, waitForActivity(m.sub) // wait for next event
	case spinner.TickMsg:
		var cmd tui.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() tui.View {
	s := fmt.Sprintf("\n %s Events received: %d\n\n Press any key to exit\n", m.spinner.View(), m.responses)
	if m.quitting {
		s += "\n"
	}
	return tui.NewView(s)
}

func main() {
	p := tui.NewProgram(model{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
