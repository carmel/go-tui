package main

import (
	"image/color"
	"log"

	"github.com/carmel/go-tui"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/ansi"
	"github.com/lucasb-eyer/go-colorful"
)

var myFancyColor color.Color

type model struct{}

var _ tui.Model = model{}

// Init implements tui.Model.
func (m model) Init() tui.Cmd {
	return tui.Batch(
		tui.RequestCapability("RGB"),
		tui.RequestCapability("Tc"),
	)
}

// Update implements tui.Model.
func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		return m, tui.Quit
	case tui.ColorProfileMsg:
		return m, tui.Println("Color profile manually set to ", msg)
	}
	return m, nil
}

// View implements tui.Model.
func (m model) View() tui.View {
	return tui.NewView("This will produce the wrong colors on Apple Terminal :)\n\n" +
		ansi.Style{}.ForegroundColor(myFancyColor).Styled("Howdy!") +
		"\n\n" +
		"Press any key to exit.")
}

func main() {
	myFancyColor, _ = colorful.Hex("#6b50ff")

	p := tui.NewProgram(model{}, tui.WithColorProfile(colorprofile.TrueColor))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
