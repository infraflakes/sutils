package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func RunTUI() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	m := &model{
		currentDir: cwd,
	}
	m.updateEntries()
	m.goUp()

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	res := finalModel.(*model)
	if res.finalPath != "" {
		return res.finalPath, nil
	}

	return "", nil
}
