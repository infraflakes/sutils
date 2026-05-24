package todo

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Initialize(configFilePath string) Model {
	var finalPath string
	if configFilePath != "" {
		absPath, err := filepath.Abs(configFilePath)
		if err != nil {
			finalPath = configFilePath
		} else {
			finalPath = absPath
		}
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil || homeDir == "" {
			finalPath = filepath.Join(".", ".cache", "sutils", "todo", "note.json")
		} else {
			finalPath = filepath.Join(homeDir, ".cache", "sutils", "todo", "note.json")
		}
	}

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 70

	dateInputs := make([]textinput.Model, 3)
	for i := range dateInputs {
		dateInputs[i] = textinput.New()
		dateInputs[i].Focus()
		dateInputs[i].CharLimit = 4
		dateInputs[i].Width = 10
	}

	m := Model{
		TextInput:      ti,
		DateInputs:     dateInputs,
		KeyMap:         DefaultKeyMap(),
		Help:           help.New(),
		ConfigFilePath: finalPath,
		MaxHistory:     50,
		ViewMode:       NormalView,
	}

	m.LoadConfig()
	m.UpdateContexts()

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
