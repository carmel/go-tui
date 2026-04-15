package main

import (
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/lipgloss/tree"
)

func main() {
	t := tree.Root(".").
		Child("macOS").
		Child(
			tree.New().
				Root("Linux").
				Child("NixOS").
				Child("Arch Linux (btw)").
				Child("Void Linux"),
		).
		Child(
			tree.New().
				Root("BSD").
				Child("FreeBSD").
				Child("OpenBSD"),
		)

	lipgloss.Println(t)
}
