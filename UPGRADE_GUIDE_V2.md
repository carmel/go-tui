# Bubble Tea v2 Upgrade Guide

This guide covers everything you need to change when upgrading from Bubble Tea v1 to v2. For a tour of all the exciting new features, check out the [What's New](https://github.com/carmel/go-tui/releases/tag/v2.0.0) doc.

> [!NOTE]
> We don't take API changes lightly and strive to make the upgrade process as simple as possible. If something feels way off, let us know.

## Migration Checklist

Here's the short version — a checklist you can follow top to bottom. Each item links to the relevant section below.

- [ ] [Update import paths](#import-paths)
- [ ] [Change `View() string` to `View() tui.View`](#view-returns-a-teaview-now)
- [ ] [Replace `tui.KeyMsg` with `tui.KeyPressMsg`](#key-messages)
- [ ] [Update key fields: `msg.Type` / `msg.Runes` / `msg.Alt`](#key-messages)
- [ ] [Replace `case " ":` with `case "space":`](#key-messages)
- [ ] [Update mouse message usage](#mouse-messages)
- [ ] [Rename mouse button constants](#mouse-messages)
- [ ] [Remove old program options → use View fields](#removed-program-options)
- [ ] [Remove imperative commands → use View fields](#removed-commands)
- [ ] [Remove old program methods](#removed-program-methods)
- [ ] [Rename `tui.WindowSize()` → `tui.RequestWindowSize`](#renamed-apis)
- [ ] [Replace `tui.Sequentially(...)` → `tui.Sequence(...)`](#renamed-apis)

## Import Paths

The module path changed to a vanity domain. Lip Gloss moved too.

```go
// Before
import "github.com/carmel/go-tui"
import "github.com/charmbracelet/lipgloss"

// After
import "github.com/carmel/go-tui"
import "github.com/carmel/go-tui/lipgloss"
```

## The Big Idea: Declarative Views

The single biggest change in v2 is the shift from **imperative commands** to **declarative View fields**. In v1, you'd use program options like `tui.WithAltScreen()` and commands like `tui.EnterAltScreen` to toggle terminal features on and off. In v2, you just set fields on the `tui.View` struct in your `View()` method and Bubble Tea handles the rest.

This means: no more startup option flags, no more toggle commands, no more fighting over state. Just declare what you want and Bubble Tea will make it so.

```go
// v1: imperative — scattered across NewProgram, Init, and Update
p := tui.NewProgram(model{}, tui.WithAltScreen(), tui.WithMouseCellMotion())

// v2: declarative — everything lives in View()
func (m model) View() tui.View {
    v := tui.NewView("Hello!")
    v.AltScreen = true
    v.MouseMode = tui.MouseModeCellMotion
    return v
}
```

Keep this in mind as you go through the rest of the guide — most of the "removed" things simply moved into View fields.

## View Returns a `tui.View` Now

The `View()` method no longer returns a `string`. It returns a `tui.View` struct.

```go
// Before:
func (m model) View() string {
    return "Hello, world!"
}

// After:
func (m model) View() tui.View {
    return tui.NewView("Hello, world!")
}
```

You can also use the longer form if you need to set additional fields:

```go
func (m model) View() tui.View {
    var v tui.View
    v.SetContent("Hello, world!")
    v.AltScreen = true
    return v
}
```

The `tui.View` struct has fields for everything that used to be controlled by options and commands:

| View Field                  | What It Does                                                    |
| --------------------------- | --------------------------------------------------------------- |
| `Content`                   | The rendered string (set via `SetContent()` or `NewView()`)     |
| `AltScreen`                 | Enter/exit the alternate screen buffer                          |
| `MouseMode`                 | `MouseModeNone`, `MouseModeCellMotion`, or `MouseModeAllMotion` |
| `ReportFocus`               | Enable focus/blur event reporting                               |
| `DisableBracketedPasteMode` | Disable bracketed paste                                         |
| `WindowTitle`               | Set the terminal window title                                   |
| `Cursor`                    | Control cursor position, shape, color, and blink                |
| `ForegroundColor`           | Set the terminal foreground color                               |
| `BackgroundColor`           | Set the terminal background color                               |
| `ProgressBar`               | Show a native terminal progress bar                             |
| `KeyboardEnhancements`      | Request keyboard enhancement features                           |
| `OnMouse`                   | Intercept mouse messages based on view content                  |

## Key Messages

Key messages got a major overhaul. Here's the quick rundown:

### `tui.KeyMsg` is now an interface

In v1, `tui.KeyMsg` was a struct you'd match on for key presses. In v2, it's an **interface** that covers both key presses and releases. For most code, you want `tui.KeyPressMsg`:

```go
// Before:
case tui.KeyMsg:
    switch msg.String() {
    case "q":
        return m, tui.Quit
    }

// After:
case tui.KeyPressMsg:
    switch msg.String() {
    case "q":
        return m, tui.Quit
    }
```

If you want to handle both presses and releases, use `tui.KeyMsg` and type-switch inside:

```go
case tui.KeyMsg:
    switch key := msg.(type) {
    case tui.KeyPressMsg:
        // key press
    case tui.KeyReleaseMsg:
        // key release
    }
```

### Key fields changed

| v1             | v2         | Notes                                                          |
| -------------- | ---------- | -------------------------------------------------------------- |
| `msg.Type`     | `msg.Code` | A `rune` — can be `tui.KeyEnter`, `'a'`, etc.                  |
| `msg.Runes`    | `msg.Text` | Now a `string`, not `[]rune`                                   |
| `msg.Alt`      | `msg.Mod`  | `msg.Mod.Contains(tui.ModAlt)` for alt, etc.                   |
| `tui.KeyRune`  | —          | Check `len(msg.Text) > 0` instead                              |
| `tui.KeyCtrlC` | —          | Use `msg.String() == "ctrl+c"` or check `msg.Code` + `msg.Mod` |

### Space bar changed

Space bar now returns `"space"` instead of `" "` when using `msg.String()`:

```go
// Before:
case " ":

// After:
case "space":
```

`key.Code` is still `' '` and `key.Text` is still `" "`, but `String()` returns `"space"`.

### Ctrl+key matching

```go
// Before:
case tui.KeyCtrlC:
    // ctrl+c

// After (option A — string matching):
case tui.KeyPressMsg:
    switch msg.String() {
    case "ctrl+c":
        // ctrl+c
    }

// After (option B — field matching):
case tui.KeyPressMsg:
    if msg.Code == 'c' && msg.Mod == tui.ModCtrl {
        // ctrl+c
    }
```

### New Key fields

These are new in v2 and don't have v1 equivalents:

- **`key.ShiftedCode`** — the shifted key code (e.g., `'B'` when pressing shift+b)
- **`key.BaseCode`** — the key on a US PC-101 layout (handy for international keyboards)
- **`key.IsRepeat`** — whether the key is auto-repeating (Kitty protocol / Windows Console only)
- **`key.Keystroke()`** — like `String()` but always includes modifier info

## Paste Messages

Paste events no longer come in as `tui.KeyMsg` with a `Paste` flag. They're now their own message types:

```go
// Before:
case tui.KeyMsg:
    if msg.Paste {
        m.text += string(msg.Runes)
    }

// After:
case tui.PasteMsg:
    m.text += msg.Content
case tui.PasteStartMsg:
    // paste started
case tui.PasteEndMsg:
    // paste ended
```

## Mouse Messages

### `tui.MouseMsg` is now an interface

In v1, `tui.MouseMsg` was a struct with `X`, `Y`, `Button`, etc. In v2, it's an **interface**. You get the coordinates by calling `msg.Mouse()`:

```go
// Before:
case tui.MouseMsg:
    x, y := msg.X, msg.Y

// After:
case tui.MouseMsg:
    mouse := msg.Mouse()
    x, y := mouse.X, mouse.Y
```

### Mouse events are split by type

Instead of checking `msg.Action`, match on specific message types:

```go
// Before:
case tui.MouseMsg:
    if msg.Action == tui.MouseActionPress && msg.Button == tui.MouseButtonLeft {
        // left click
    }

// After:
case tui.MouseClickMsg:
    if msg.Button == tui.MouseLeft {
        // left click
    }
case tui.MouseReleaseMsg:
    // release
case tui.MouseWheelMsg:
    // scroll
case tui.MouseMotionMsg:
    // movement
```

### Button constants renamed

| v1                          | v2                    |
| --------------------------- | --------------------- |
| `tui.MouseButtonLeft`       | `tui.MouseLeft`       |
| `tui.MouseButtonRight`      | `tui.MouseRight`      |
| `tui.MouseButtonMiddle`     | `tui.MouseMiddle`     |
| `tui.MouseButtonWheelUp`    | `tui.MouseWheelUp`    |
| `tui.MouseButtonWheelDown`  | `tui.MouseWheelDown`  |
| `tui.MouseButtonWheelLeft`  | `tui.MouseWheelLeft`  |
| `tui.MouseButtonWheelRight` | `tui.MouseWheelRight` |

### `tui.MouseEvent` → `tui.Mouse`

The `MouseEvent` struct is gone. The new `Mouse` struct has `X`, `Y`, `Button`, and `Mod` fields.

### Mouse mode is now a View field

```go
// Before:
p := tui.NewProgram(model{}, tui.WithMouseCellMotion())

// After:
func (m model) View() tui.View {
    v := tui.NewView("...")
    v.MouseMode = tui.MouseModeCellMotion
    return v
}
```

## Removed Program Options

These options no longer exist. They all moved to View fields.

| Removed Option                | Do This Instead                                                      |
| ----------------------------- | -------------------------------------------------------------------- |
| `tui.WithAltScreen()`         | `view.AltScreen = true`                                              |
| `tui.WithMouseCellMotion()`   | `view.MouseMode = tui.MouseModeCellMotion`                           |
| `tui.WithMouseAllMotion()`    | `view.MouseMode = tui.MouseModeAllMotion`                            |
| `tui.WithReportFocus()`       | `view.ReportFocus = true`                                            |
| `tui.WithoutBracketedPaste()` | `view.DisableBracketedPasteMode = true`                              |
| `tui.WithInputTTY()`          | Just remove it — v2 always opens the TTY for input automatically     |
| `tui.WithANSICompressor()`    | Just remove it — the new renderer handles optimization automatically |

## Removed Commands

These commands no longer exist. Set the corresponding View field instead.

| Removed Command             | Do This Instead                                           |
| --------------------------- | --------------------------------------------------------- |
| `tui.EnterAltScreen`        | `view.AltScreen = true`                                   |
| `tui.ExitAltScreen`         | `view.AltScreen = false`                                  |
| `tui.EnableMouseCellMotion` | `view.MouseMode = tui.MouseModeCellMotion`                |
| `tui.EnableMouseAllMotion`  | `view.MouseMode = tui.MouseModeAllMotion`                 |
| `tui.DisableMouse`          | `view.MouseMode = tui.MouseModeNone`                      |
| `tui.HideCursor`            | `view.Cursor = nil`                                       |
| `tui.ShowCursor`            | `view.Cursor = &tui.Cursor{...}` or `tui.NewCursor(x, y)` |
| `tui.EnableBracketedPaste`  | `view.DisableBracketedPasteMode = false`                  |
| `tui.DisableBracketedPaste` | `view.DisableBracketedPasteMode = true`                   |
| `tui.EnableReportFocus`     | `view.ReportFocus = true`                                 |
| `tui.DisableReportFocus`    | `view.ReportFocus = false`                                |
| `tui.SetWindowTitle("...")` | `view.WindowTitle = "..."`                                |

## Removed Program Methods

These methods on `*Program` are gone.

| Removed Method               | Do This Instead                                  |
| ---------------------------- | ------------------------------------------------ |
| `p.Start()`                  | `p.Run()`                                        |
| `p.StartReturningModel()`    | `p.Run()`                                        |
| `p.EnterAltScreen()`         | `view.AltScreen = true` in `View()`              |
| `p.ExitAltScreen()`          | `view.AltScreen = false` in `View()`             |
| `p.EnableMouseCellMotion()`  | `view.MouseMode` in `View()`                     |
| `p.DisableMouseCellMotion()` | `view.MouseMode = tui.MouseModeNone` in `View()` |
| `p.EnableMouseAllMotion()`   | `view.MouseMode` in `View()`                     |
| `p.DisableMouseAllMotion()`  | `view.MouseMode = tui.MouseModeNone` in `View()` |
| `p.SetWindowTitle(...)`      | `view.WindowTitle` in `View()`                   |

## Renamed APIs

| v1                      | v2                      | Notes                                       |
| ----------------------- | ----------------------- | ------------------------------------------- |
| `tui.Sequentially(...)` | `tui.Sequence(...)`     | `Sequentially` was already deprecated in v1 |
| `tui.WindowSize()`      | `tui.RequestWindowSize` | Now returns `Msg` directly, not a `Cmd`     |

## New Program Options

These are new in v2:

| Option                     | What It Does                                       |
| -------------------------- | -------------------------------------------------- |
| `tui.WithColorProfile(p)`  | Force a specific color profile (great for testing) |
| `tui.WithWindowSize(w, h)` | Set initial terminal size (great for testing)      |

## Complete Before & After

Here's a minimal but complete program showing the most common migration patterns side by side.

**v1:**

```go
package main

import (
    "fmt"
    "os"

    "github.com/carmel/go-tui"
)

type model struct {
    count int
}

func (m model) Init() tui.Cmd {
    return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
    switch msg := msg.(type) {
    case tui.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tui.Quit
        case " ":
            m.count++
        }
    case tui.MouseMsg:
        if msg.Action == tui.MouseActionPress && msg.Button == tui.MouseButtonLeft {
            m.count++
        }
    }
    return m, nil
}

func (m model) View() string {
    return fmt.Sprintf("Count: %d\n\nSpace or click to increment. q to quit.\n", m.count)
}

func main() {
    p := tui.NewProgram(model{}, tui.WithAltScreen(), tui.WithMouseCellMotion())
    if _, err := p.Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**v2:**

```go
package main

import (
    "fmt"
    "os"

    "github.com/carmel/go-tui"
)

type model struct {
    count int
}

func (m model) Init() tui.Cmd {
    return nil
}

func (m model) Update(msg tui.Msg) (tui.Model, tui.Cmd) {
    switch msg := msg.(type) {
    case tui.KeyPressMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tui.Quit
        case "space":
            m.count++
        }
    case tui.MouseClickMsg:
        if msg.Button == tui.MouseLeft {
            m.count++
        }
    }
    return m, nil
}

func (m model) View() tui.View {
    v := tui.NewView(fmt.Sprintf("Count: %d\n\nSpace or click to increment. q to quit.\n", m.count))
    v.AltScreen = true
    v.MouseMode = tui.MouseModeCellMotion
    return v
}

func main() {
    p := tui.NewProgram(model{})
    if _, err := p.Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

Notice how the `NewProgram` call got simpler? All the terminal feature flags moved into `View()` where they belong.

## Quick Reference

A flat old → new lookup table. Handy for search-and-replace and LLM-assisted migration.

### Import Paths

| v1                                  | v2                                  |
| ----------------------------------- | ----------------------------------- |
| `github.com/carmel/go-tui`          | `github.com/carmel/go-tui`          |
| `github.com/charmbracelet/lipgloss` | `github.com/carmel/go-tui/lipgloss` |

### Model Interface

| v1              | v2                |
| --------------- | ----------------- |
| `View() string` | `View() tui.View` |

### Key Events

| v1                    | v2                                                                        |
| --------------------- | ------------------------------------------------------------------------- |
| `tui.KeyMsg` (struct) | `tui.KeyPressMsg` for presses, `tui.KeyMsg` (interface) for both          |
| `msg.Type`            | `msg.Code`                                                                |
| `msg.Runes`           | `msg.Text` (string, not `[]rune`)                                         |
| `msg.Alt`             | `msg.Mod.Contains(tui.ModAlt)`                                            |
| `tui.KeyRune`         | check `len(msg.Text) > 0`                                                 |
| `tui.KeyCtrlC`        | `msg.Code == 'c' && msg.Mod == tui.ModCtrl` or `msg.String() == "ctrl+c"` |
| `case " ":` (space)   | `case "space":`                                                           |

### Mouse Events

| v1                         | v2                                                        |
| -------------------------- | --------------------------------------------------------- |
| `tui.MouseMsg` (struct)    | `tui.MouseMsg` (interface) — call `.Mouse()` for the data |
| `tui.MouseEvent`           | `tui.Mouse`                                               |
| `tui.MouseButtonLeft`      | `tui.MouseLeft`                                           |
| `tui.MouseButtonRight`     | `tui.MouseRight`                                          |
| `tui.MouseButtonMiddle`    | `tui.MouseMiddle`                                         |
| `tui.MouseButtonWheelUp`   | `tui.MouseWheelUp`                                        |
| `tui.MouseButtonWheelDown` | `tui.MouseWheelDown`                                      |
| `msg.X`, `msg.Y` (direct)  | `msg.Mouse().X`, `msg.Mouse().Y`                          |

### Options → View Fields

| v1 Option                     | v2 View Field                              |
| ----------------------------- | ------------------------------------------ |
| `tui.WithAltScreen()`         | `view.AltScreen = true`                    |
| `tui.WithMouseCellMotion()`   | `view.MouseMode = tui.MouseModeCellMotion` |
| `tui.WithMouseAllMotion()`    | `view.MouseMode = tui.MouseModeAllMotion`  |
| `tui.WithReportFocus()`       | `view.ReportFocus = true`                  |
| `tui.WithoutBracketedPaste()` | `view.DisableBracketedPasteMode = true`    |

### Commands → View Fields

| v1 Command                                               | v2 View Field                                          |
| -------------------------------------------------------- | ------------------------------------------------------ |
| `tui.EnterAltScreen` / `tui.ExitAltScreen`               | `view.AltScreen = true/false`                          |
| `tui.EnableMouseCellMotion`                              | `view.MouseMode = tui.MouseModeCellMotion`             |
| `tui.EnableMouseAllMotion`                               | `view.MouseMode = tui.MouseModeAllMotion`              |
| `tui.DisableMouse`                                       | `view.MouseMode = tui.MouseModeNone`                   |
| `tui.HideCursor` / `tui.ShowCursor`                      | `view.Cursor = nil` / `view.Cursor = &tui.Cursor{...}` |
| `tui.EnableBracketedPaste` / `tui.DisableBracketedPaste` | `view.DisableBracketedPasteMode = false/true`          |
| `tui.EnableReportFocus` / `tui.DisableReportFocus`       | `view.ReportFocus = true/false`                        |
| `tui.SetWindowTitle("...")`                              | `view.WindowTitle = "..."`                             |

### Removed Options (No Replacement Needed)

| v1 Option                  | What Happened                                       |
| -------------------------- | --------------------------------------------------- |
| `tui.WithInputTTY()`       | v2 always opens the TTY for input automatically     |
| `tui.WithANSICompressor()` | The new renderer handles optimization automatically |

### Removed Program Methods

| v1 Method                    | v2 Replacement                                   |
| ---------------------------- | ------------------------------------------------ |
| `p.Start()`                  | `p.Run()`                                        |
| `p.StartReturningModel()`    | `p.Run()`                                        |
| `p.EnterAltScreen()`         | `view.AltScreen = true` in `View()`              |
| `p.ExitAltScreen()`          | `view.AltScreen = false` in `View()`             |
| `p.EnableMouseCellMotion()`  | `view.MouseMode` in `View()`                     |
| `p.DisableMouseCellMotion()` | `view.MouseMode = tui.MouseModeNone` in `View()` |
| `p.EnableMouseAllMotion()`   | `view.MouseMode` in `View()`                     |
| `p.DisableMouseAllMotion()`  | `view.MouseMode = tui.MouseModeNone` in `View()` |
| `p.SetWindowTitle(...)`      | `view.WindowTitle` in `View()`                   |

### Other Renames

| v1                      | v2                                                     |
| ----------------------- | ------------------------------------------------------ |
| `tui.Sequentially(...)` | `tui.Sequence(...)`                                    |
| `tui.WindowSize()`      | `tui.RequestWindowSize` (now returns `Msg`, not `Cmd`) |

### New Program Options

| Option                     | Description                                 |
| -------------------------- | ------------------------------------------- |
| `tui.WithColorProfile(p)`  | Force a specific color profile              |
| `tui.WithWindowSize(w, h)` | Set initial window size (great for testing) |

## Feedback

Have thoughts on the v2 upgrade? We'd _love_ to hear about it. Let us know on…

- [Discord](https://charm.land/chat)
- [Matrix](https://charm.land/matrix)
- [Email](mailto:vt100@charm.land)

---

Part of [Charm](https://charm.land).

<a href="https://charm.land/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source • نحنُ نحب المصادر المفتوحة
