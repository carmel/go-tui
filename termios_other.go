//go:build !windows && !darwin && !dragonfly && !freebsd && !linux && !solaris && !aix
// +build !windows,!darwin,!dragonfly,!freebsd,!linux,!solaris,!aix

package tui

import "github.com/charmbracelet/x/term"

func (*Program) checkOptimizedMovements(*term.State) {}
