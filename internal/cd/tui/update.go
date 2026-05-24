package tui

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.searchMode {
		return m.updateSearch(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "/":
			m.searchMode = true
			m.searchQuery = ""

		case ".":
			m.showHidden = !m.showHidden
			m.updateEntries()

		case "backspace":
			m.showFiles = !m.showFiles
			m.updateEntries()

		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
				m.updateEntries()
			}

		case "down", "j":
			if m.selectedIdx < len(m.currentEntries)-1 {
				m.selectedIdx++
				m.updateEntries()
			}

		case "left", "h":
			m.goUp()

		case "right", "l":
			m.goIn()

		case "enter":
			if len(m.currentEntries) > 0 {
				sel := m.currentEntries[m.selectedIdx]
				if sel.isDir {
					m.finalPath = sel.path
					return m, tea.Quit
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.searchMode = false
			m.searchQuery = ""
			return m, nil
		case tea.KeyEnter:
			m.searchMode = false
			return m, nil
		case tea.KeyBackspace:
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.jumpToMatch()
			} else {
				m.searchMode = false
			}
		default:
			if msg.Type == tea.KeyRunes {
				m.searchQuery += string(msg.Runes)
				m.jumpToMatch()
			}
		}
	}
	return m, nil
}

func (m *model) jumpToMatch() {
	if m.searchQuery == "" {
		return
	}
	for i, e := range m.currentEntries {
		if strings.Contains(strings.ToLower(e.name), strings.ToLower(m.searchQuery)) {
			m.selectedIdx = i
			m.updateEntries()
			break
		}
	}
}

func (m *model) goUp() {
	parent := filepath.Dir(m.currentDir)
	if parent != m.currentDir {
		oldDir := filepath.Base(m.currentDir)
		m.currentDir = parent
		m.updateEntries()
		for i, e := range m.currentEntries {
			if e.name == oldDir {
				m.selectedIdx = i
				break
			}
		}
		m.updateEntries()
	}
}

func (m *model) goIn() {
	if len(m.currentEntries) > 0 {
		sel := m.currentEntries[m.selectedIdx]
		if sel.isDir {
			m.currentDir = sel.path
			m.selectedIdx = 0
			m.updateEntries()
		}
	}
}
