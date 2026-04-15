package main

import (
	"github.com/carmel/go-tui/lipgloss"
	"github.com/carmel/go-tui/lipgloss/list"
)

func main() {
	l := list.New(
		"A",
		"B",
		"C",
		list.New(
			"D",
			"E",
			"F",
		).Enumerator(list.Roman),
		"G",
	)
	lipgloss.Println(l)
}
