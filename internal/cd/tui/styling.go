package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	Renderer = lipgloss.NewRenderer(os.Stderr)

	// HeaderStyle uses inverted text (background <-> foreground) + bolding
	HeaderStyle = Renderer.NewStyle().
			Bold(true).
			Reverse(true).
			Padding(0, 1)

	// SelectedStyle safely highlights the active item by flipping text/bg colors
	SelectedStyle = Renderer.NewStyle().
			Bold(true).
			Reverse(true)

	// DimStyle drops the intensity slightly if the terminal supports faint,
	// otherwise defaults gracefully to standard foreground
	DimStyle = Renderer.NewStyle().
			Faint(true)

	// BrightStyle uses standard bold text for clean contrast without color clashing
	BrightStyle = Renderer.NewStyle().
			Bold(true)

	// KeyStyle uses raw underline or bolding to emphasize shortcuts rather than breaking the mono aesthetic
	KeyStyle = Renderer.NewStyle().
			Bold(true).
			Underline(true)

	// MatchStyle handles fuzzy match highlighting using standard bolding, or an absolute contrast flip
	MatchStyle = Renderer.NewStyle().
			Bold(true).
			Underline(true)

	// Raw ANSI overrides for in-line string manipulations
	// \x1b[7m turns on Inverse/Reverse video mode, \x1b[27m turns it off
	MatchStyleInSelected = "\x1b[7m"
	RestoreSelectedBg    = "\x1b[27m"

	// Borders use clean, native terminal separators instead of color indexing
	ColumnStyle = Renderer.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			Padding(0, 1)

	CurrentColumnStyle = Renderer.NewStyle().
				Border(lipgloss.DoubleBorder(), false, true, false, true).
				Padding(0, 1)
)
