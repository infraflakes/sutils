package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	Renderer = lipgloss.NewRenderer(os.Stderr)

	HeaderStyle = Renderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("4")).
			Padding(0, 1)

	SelectedStyle = Renderer.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("4")).
			Bold(true)

	DimStyle = Renderer.NewStyle().
			Foreground(lipgloss.Color("8"))

	BrightStyle = Renderer.NewStyle().
			Foreground(lipgloss.Color("7"))

	KeyStyle = Renderer.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	MatchStyle = Renderer.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("11")).
			Bold(true)

	MatchStyleInSelected  = "\x1b[103m"
	RestoreSelectedBg     = "\x1b[44m"
	ColumnStyle           = Renderer.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, false).BorderForeground(lipgloss.Color("8")).Padding(0, 1)
	CurrentColumnStyle    = Renderer.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, true).BorderForeground(lipgloss.Color("4")).Padding(0, 1)
)
