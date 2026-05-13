package todo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		m.Help.Width = msg.Width
		m.TextInput.Width = msg.Width - 20
		return m, tea.ClearScreen

	case tea.KeyMsg:
		if m.HelpVisible {
			switch {
			case key.Matches(msg, m.KeyMap.Help), key.Matches(msg, m.KeyMap.Back):
				m.HelpVisible = false
			}
			return m, nil
		}

		m.ErrorMessage = ""

		switch m.ViewMode {
		case InputView:
			return m.UpdateInputMode(msg)
		case DateInputView:
			return m.UpdateDateInputMode(msg)
		case RemoveTagView:
			return m.UpdateRemoveTagMode(msg)
		}

		switch m.ViewMode {
		case NormalView:
			return m.UpdateNormalView(msg)
		case KanbanView:
			return m.UpdateKanbanView(msg)
		case StatsView:
			return m.UpdateStatsView(msg)
		}
	}

	return m, nil
}

func (m Model) UpdateInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.KeyMap.Back):
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Enter):
		input := strings.TrimSpace(m.TextInput.Value())
		m.TextInput.SetValue("")

		switch m.InputMode {
		case AddTaskInput:
			if input != "" {
				m.SaveStateForUndo()
				m.AddTask(input)
				m.SaveConfig()
			}
		case EditTaskInput:
			if input != "" {
				m.SaveStateForUndo()
				m.EditCurrentTask(input)
				m.SaveConfig()
			}
		case AddContextInput:
			if input != "" {
				m.AddContext(input)
				m.SaveConfig()
			}
		case RenameContextInput:
			if input != "" && input != m.CurrentContext {
				m.RenameContext(input)
				m.SaveConfig()
			}
		case AddTagInput:
			if input != "" {
				m.SaveStateForUndo()
				m.AddTagToCurrentTask(input)
				m.SaveConfig()
			}
		case DeleteConfirmInput:
			if strings.ToLower(input) == "y" {
				m.SaveStateForUndo()
				m.DeleteContext()
				m.SaveConfig()
			}
		}

		m.ViewMode = NormalView
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) UpdateDateInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.KeyMap.Back):
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Enter):
		day := m.DateInputs[0].Value()
		month := m.DateInputs[1].Value()
		year := m.DateInputs[2].Value()
		dateStr := fmt.Sprintf("%s-%s-%s", year, month, day)
		m.SaveStateForUndo()
		m.SetDueDateForCurrentTask(dateStr)
		m.SaveConfig()
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Up):
		m.DateInputs[m.DateInputIndex].Blur()
		m.DateInputIndex = (m.DateInputIndex - 1 + 3) % 3
		m.DateInputs[m.DateInputIndex].Focus()

	case key.Matches(msg, m.KeyMap.Down):
		m.DateInputs[m.DateInputIndex].Blur()
		m.DateInputIndex = (m.DateInputIndex + 1) % 3
		m.DateInputs[m.DateInputIndex].Focus()
	}

	m.DateInputs[m.DateInputIndex], cmd = m.DateInputs[m.DateInputIndex].Update(msg)
	return m, cmd
}

func (m Model) UpdateRemoveTagMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Back):
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Enter):
		m.SaveStateForUndo()
		m.RemoveTagsFromCurrentTask()
		m.SaveConfig()
		m.ViewMode = NormalView
		return m, nil

	case key.Matches(msg, m.KeyMap.Up):
		if m.RemoveTagIndex > 0 {
			m.RemoveTagIndex--
		}

	case key.Matches(msg, m.KeyMap.Down):
		task := m.GetCurrentTask()
		if m.RemoveTagIndex < len(task.Tags)-1 {
			m.RemoveTagIndex++
		}

	case key.Matches(msg, m.KeyMap.Toggle):
		m.RemoveTagChecks[m.RemoveTagIndex] = !m.RemoveTagChecks[m.RemoveTagIndex]
	}

	return m, nil
}

func (m Model) UpdateNormalView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.KeyMap.Back):
		return m, nil

	case key.Matches(msg, m.KeyMap.Up):
		if m.MovingMode {
			m.MoveTaskUp()
		} else {
			m.MoveUp()
		}

	case key.Matches(msg, m.KeyMap.Down):
		if m.MovingMode {
			m.MoveTaskDown()
		} else {
			m.MoveDown()
		}

	case key.Matches(msg, m.KeyMap.Left):
		m.PreviousContext()

	case key.Matches(msg, m.KeyMap.Right):
		m.NextContext()

	case key.Matches(msg, m.KeyMap.Toggle):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.ToggleCurrentTask()
			m.SaveConfig()
		}

	case key.Matches(msg, m.KeyMap.Add):
		m.ShowInputDialog(AddTaskInput, "Add new task:")

	case key.Matches(msg, m.KeyMap.Edit):
		if len(m.GetFilteredTasks()) > 0 {
			task := m.GetCurrentTask()
			m.ShowInputDialog(EditTaskInput, "Edit task:")
			m.TextInput.SetValue(task.Task)
		}

	case key.Matches(msg, m.KeyMap.Delete):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.DeleteCurrentTask()
			m.SaveConfig()
		}

	case key.Matches(msg, m.KeyMap.AddContext):
		m.ShowInputDialog(AddContextInput, "New context name:")

	case key.Matches(msg, m.KeyMap.RenameContext):
		m.ShowInputDialog(RenameContextInput, "Rename context to:")
		m.TextInput.SetValue(m.CurrentContext)

	case key.Matches(msg, m.KeyMap.DeleteContext):
		if len(m.Contexts) > 1 {
			m.ShowInputDialog(DeleteConfirmInput, fmt.Sprintf("Delete context '%s'? (y/n):", m.CurrentContext))
		} else {
			m.ErrorMessage = "Cannot delete the only context"
		}

	case key.Matches(msg, m.KeyMap.TogglePriority):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.ToggleCurrentTaskPriority()
			m.SaveConfig()
		}

	case key.Matches(msg, m.KeyMap.AddTag):
		if len(m.GetFilteredTasks()) > 0 {
			m.ShowInputDialog(AddTagInput, "Add tag:")
		}

	case key.Matches(msg, m.KeyMap.RemoveTag):
		if len(m.GetFilteredTasks()) > 0 {
			m.ShowRemoveTagDialog()
		}

	case key.Matches(msg, m.KeyMap.SetDueDate):
		if len(m.GetFilteredTasks()) > 0 {
			m.ShowDateInputDialog()
		}

	case key.Matches(msg, m.KeyMap.ClearDueDate):
		if len(m.GetFilteredTasks()) > 0 {
			m.SaveStateForUndo()
			m.SetDueDateForCurrentTask("clear")
			m.SaveConfig()
		}

	case key.Matches(msg, m.KeyMap.KanbanView):
		m.ViewMode = KanbanView

	case key.Matches(msg, m.KeyMap.StatsView):
		m.ViewMode = StatsView

	case key.Matches(msg, m.KeyMap.Undo):
		m.Undo()
		m.SaveConfig()

	case key.Matches(msg, m.KeyMap.Help):
		m.HelpVisible = true

	case key.Matches(msg, m.KeyMap.Move):
		if len(m.GetFilteredTasks()) > 0 {
			m.MovingMode = !m.MovingMode
			if m.MovingMode {
				m.MovingTaskID = m.GetCurrentTask().ID
			} else {
				m.SaveStateForUndo()
				m.SaveConfig()
			}
		}
	}

	return m, nil
}

func (m Model) UpdateKanbanView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Back), key.Matches(msg, m.KeyMap.Quit), key.Matches(msg, m.KeyMap.KanbanView):
		m.ViewMode = NormalView
		m.KanbanScrollY = 0
		m.KanbanScrollX = 0
	case key.Matches(msg, m.KeyMap.Up):
		if m.KanbanScrollY > 0 {
			m.KanbanScrollY--
		}
	case key.Matches(msg, m.KeyMap.Down):
		m.KanbanScrollY++
	case key.Matches(msg, m.KeyMap.Left):
		if m.KanbanScrollX > 0 {
			m.KanbanScrollX--
		}
	case key.Matches(msg, m.KeyMap.Right):
		m.KanbanScrollX++
	}
	return m, nil
}

func (m Model) UpdateStatsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Back), key.Matches(msg, m.KeyMap.Quit), key.Matches(msg, m.KeyMap.StatsView):
		m.ViewMode = NormalView
	}
	return m, nil
}
