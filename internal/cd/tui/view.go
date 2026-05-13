package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) View() string {
	if m.quitting {
		return ""
	}

	header := HeaderStyle.Render(m.currentDir)

	colWidth := (m.width - 6) / 3

	parentCol := ColumnStyle.Width(colWidth).Height(m.height - 4).Render(m.renderColumn(m.parentEntries, -1, false, colWidth-2))
	currentCol := CurrentColumnStyle.Width(colWidth).Height(m.height - 4).Render(m.renderColumn(m.currentEntries, m.selectedIdx, true, colWidth-2))
	previewCol := Renderer.NewStyle().Width(colWidth).Height(m.height - 4).Render(m.renderColumn(m.previewEntries, -1, false, colWidth-2))

	body := lipgloss.JoinHorizontal(lipgloss.Top, parentCol, currentCol, previewCol)

	fileStatus := " [DIRS]"
	if m.showFiles {
		fileStatus = " [ALL]"
	}
	dotStatus := ""
	if m.showHidden {
		dotStatus = " [DOTS]"
	}

	help := lipgloss.JoinHorizontal(lipgloss.Top,
		KeyStyle.Render(" ."), BrightStyle.Render(" toggle hidden "),
		KeyStyle.Render(" ⌫"), BrightStyle.Render(" toggle files "),
		KeyStyle.Render(" /"), BrightStyle.Render(" search "),
		KeyStyle.Render(" q"), BrightStyle.Render(" quit"),
	)

	if m.searchMode {
		help = lipgloss.JoinHorizontal(lipgloss.Top,
			BrightStyle.Render("Search: "), BrightStyle.Render(m.searchQuery),
		)
	}
	status := BrightStyle.Render(fmt.Sprintf("%d/%d%s%s ", m.selectedIdx+1, len(m.currentEntries), fileStatus, dotStatus))

	gapWidth := m.width - lipgloss.Width(help) - lipgloss.Width(status)
	if gapWidth < 0 {
		gapWidth = 0
	}
	gap := lipgloss.NewStyle().Width(gapWidth).Render("")

	footer := "\n" + lipgloss.JoinHorizontal(lipgloss.Top, help, gap, status)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m *model) renderColumn(entries []entry, selectedIdx int, isActive bool, width int) string {
	var s string
	height := m.height - 4

	start := 0
	if selectedIdx >= height {
		start = selectedIdx - height + 1
	}

	for i := start; i < len(entries) && i < start+height; i++ {
		e := entries[i]
		name := e.name
		if e.isDir {
			name += "/"
		}

		if i == selectedIdx && isActive {
			displayName := m.highlightMatchInSelected(name)
			padding := width - lipgloss.Width(" "+name)
			if padding < 0 {
				padding = 0
			}
			s += SelectedStyle.Render(" "+displayName+strings.Repeat(" ", padding)) + "\n"
		} else if i == selectedIdx && !isActive {
			s += DimStyle.Render("  "+m.highlightMatch(name)) + "\n"
		} else {
			s += "  " + m.highlightMatch(name) + "\n"
		}
	}
	return s
}

func (m *model) highlightMatchInSelected(name string) string {
	if m.searchQuery == "" {
		return name
	}
	idx := strings.Index(strings.ToLower(name), strings.ToLower(m.searchQuery))
	if idx == -1 {
		return name
	}
	return name[:idx] + MatchStyleInSelected + name[idx:idx+len(m.searchQuery)] + RestoreSelectedBg + name[idx+len(m.searchQuery):]
}

func (m *model) highlightMatch(name string) string {
	if m.searchQuery == "" {
		return name
	}
	idx := strings.Index(strings.ToLower(name), strings.ToLower(m.searchQuery))
	if idx == -1 {
		return name
	}
	return name[:idx] + MatchStyle.Render(name[idx:idx+len(m.searchQuery)]) + name[idx+len(m.searchQuery):]
}
