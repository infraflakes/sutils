package tui

import (
	"os"
	"path/filepath"
	"strings"
)

func (m *model) updateEntries() {
	m.currentEntries = listEntries(m.currentDir, m.showFiles, m.showHidden)
	if m.selectedIdx >= len(m.currentEntries) {
		m.selectedIdx = 0
	}
	if len(m.currentEntries) == 0 {
		m.selectedIdx = 0
	}

	parentDir := filepath.Dir(m.currentDir)
	if parentDir == m.currentDir {
		m.parentEntries = nil
	} else {
		m.parentEntries = listEntries(parentDir, m.showFiles, m.showHidden)
	}

	if len(m.currentEntries) > 0 {
		sel := m.currentEntries[m.selectedIdx]
		if sel.isDir {
			m.previewEntries = listEntries(sel.path, m.showFiles, m.showHidden)
		} else {
			m.previewEntries = nil
		}
	} else {
		m.previewEntries = nil
	}
}

func listEntries(path string, showFiles, showHidden bool) []entry {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	var entries []entry
	for _, f := range files {
		name := f.Name()

		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		if !showFiles && !f.IsDir() {
			continue
		}
		entries = append(entries, entry{
			name:  name,
			isDir: f.IsDir(),
			path:  filepath.Join(path, name),
		})
	}
	return entries
}
