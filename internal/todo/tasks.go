package todo

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

func (m *Model) ShowInputDialog(mode InputMode, prompt string) {
	m.ViewMode = InputView
	m.InputMode = mode
	m.InputPrompt = prompt
	m.TextInput.SetValue("")
	m.TextInput.Focus()
}

func (m *Model) ShowDateInputDialog() {
	m.ViewMode = DateInputView
	m.DateInputIndex = 0

	now := time.Now()
	m.DateInputs[0].SetValue(fmt.Sprintf("%02d", now.Day()))
	m.DateInputs[1].SetValue(fmt.Sprintf("%02d", now.Month()))
	m.DateInputs[2].SetValue(fmt.Sprintf("%d", now.Year()))

	for i := range m.DateInputs {
		m.DateInputs[i].Focus()
	}
}

func (m *Model) ShowRemoveTagDialog() {
	task := m.GetCurrentTask()
	if len(task.Tags) == 0 {
		m.ErrorMessage = "No tags to remove"
		return
	}

	m.ViewMode = RemoveTagView
	m.RemoveTagIndex = 0
	m.RemoveTagChecks = make([]bool, len(task.Tags))
}

func (m *Model) GetFilteredTasks() []Task {
	return m.GetTasksForContext(m.CurrentContext)
}

func (m *Model) GetTasksForContext(context string) []Task {
	var filtered []Task
	for _, task := range m.Tasks {
		if task.Context == context {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func (m *Model) GetCurrentTask() Task {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 || m.SelectedIndex >= len(tasks) {
		return Task{}
	}
	return tasks[m.SelectedIndex]
}

func (m *Model) MoveUp() {
	tasks := m.GetFilteredTasks()
	if len(tasks) > 0 {
		m.SelectedIndex = (m.SelectedIndex - 1 + len(tasks)) % len(tasks)
	}
}

func (m *Model) MoveDown() {
	tasks := m.GetFilteredTasks()
	if len(tasks) > 0 {
		m.SelectedIndex = (m.SelectedIndex + 1) % len(tasks)
	}
}

func (m *Model) findTaskIndexByID(id int) int {
	return slices.IndexFunc(m.Tasks, func(t Task) bool {
		return t.ID == id
	})
}

func (m *Model) MoveTaskUp() {
	tasks := m.GetFilteredTasks()
	if m.SelectedIndex > 0 {
		taskToMove := tasks[m.SelectedIndex]
		taskToSwapWith := tasks[m.SelectedIndex-1]
		idxMove := m.findTaskIndexByID(taskToMove.ID)
		idxSwap := m.findTaskIndexByID(taskToSwapWith.ID)
		if idxMove != -1 && idxSwap != -1 {
			m.Tasks[idxMove], m.Tasks[idxSwap] = m.Tasks[idxSwap], m.Tasks[idxMove]
			m.SelectedIndex--
		}
	}
}

func (m *Model) MoveTaskDown() {
	tasks := m.GetFilteredTasks()
	if m.SelectedIndex < len(tasks)-1 {
		taskToMove := tasks[m.SelectedIndex]
		taskToSwapWith := tasks[m.SelectedIndex+1]
		idxMove := m.findTaskIndexByID(taskToMove.ID)
		idxSwap := m.findTaskIndexByID(taskToSwapWith.ID)
		if idxMove != -1 && idxSwap != -1 {
			m.Tasks[idxMove], m.Tasks[idxSwap] = m.Tasks[idxSwap], m.Tasks[idxMove]
			m.SelectedIndex++
		}
	}
}

func (m *Model) NextContext() {
	if len(m.Contexts) > 0 {
		currentIdx := slices.Index(m.Contexts, m.CurrentContext)
		if currentIdx == -1 {
			currentIdx = 0
		}
		nextIdx := (currentIdx + 1) % len(m.Contexts)
		m.CurrentContext = m.Contexts[nextIdx]
		m.SelectedIndex = 0
	}
}

func (m *Model) PreviousContext() {
	if len(m.Contexts) > 0 {
		currentIdx := slices.Index(m.Contexts, m.CurrentContext)
		if currentIdx == -1 {
			currentIdx = 0
		}
		prevIdx := (currentIdx - 1 + len(m.Contexts)) % len(m.Contexts)
		m.CurrentContext = m.Contexts[prevIdx]
		m.SelectedIndex = 0
	}
}

func (m *Model) AddContext(contextName string) {
	if slices.Contains(m.Contexts, contextName) {
		m.ErrorMessage = "Context already exists"
		return
	}
	m.Contexts = append(m.Contexts, contextName)
	m.CurrentContext = contextName
	m.SelectedIndex = 0
}

func (m *Model) RenameContext(newName string) {
	if newName == m.CurrentContext {
		return
	}
	if slices.Contains(m.Contexts, newName) {
		m.ErrorMessage = "Context name already exists"
		return
	}
	oldName := m.CurrentContext
	if idx := slices.Index(m.Contexts, oldName); idx != -1 {
		m.Contexts[idx] = newName
	}
	for i := range m.Tasks {
		if m.Tasks[i].Context == oldName {
			m.Tasks[i].Context = newName
		}
	}
	m.CurrentContext = newName
}

func (m *Model) DeleteContext() {
	if len(m.Contexts) <= 1 {
		m.ErrorMessage = "Cannot delete the only context"
		return
	}
	m.Tasks = slices.DeleteFunc(m.Tasks, func(t Task) bool {
		return t.Context == m.CurrentContext
	})
	if idx := slices.Index(m.Contexts, m.CurrentContext); idx != -1 {
		m.Contexts = slices.Delete(m.Contexts, idx, idx+1)
	}
	if len(m.Contexts) > 0 {
		m.CurrentContext = m.Contexts[0]
		m.SelectedIndex = 0
	}
}

func (m *Model) UpdateContexts() {
	uniqueContexts := make(map[string]bool)
	for _, task := range m.Tasks {
		uniqueContexts[task.Context] = true
	}
	for _, ctx := range m.Contexts {
		uniqueContexts[ctx] = true
	}
	m.Contexts = make([]string, 0, len(uniqueContexts))
	for context := range uniqueContexts {
		m.Contexts = append(m.Contexts, context)
	}
	slices.Sort(m.Contexts)
	if m.CurrentContext == "" || !slices.Contains(m.Contexts, m.CurrentContext) {
		if len(m.Contexts) > 0 {
			m.CurrentContext = m.Contexts[0]
		} else {
			m.CurrentContext = "Work"
			m.Contexts = []string{"Work"}
		}
	}
}

func (m *Model) ToggleCurrentTask() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}
	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		m.Tasks[idx].Checked = !m.Tasks[idx].Checked
	}
}

func (m *Model) AddTask(taskText string) {
	newTask := Task{
		ID:      m.NextID,
		Task:    taskText,
		Checked: false,
		Context: m.CurrentContext,
	}
	m.Tasks = append(m.Tasks, newTask)
	m.NextID++
	filtered := m.GetFilteredTasks()
	m.SelectedIndex = len(filtered) - 1
}

func (m *Model) EditCurrentTask(newText string) {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}
	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		m.Tasks[idx].Task = newText
	}
}

func (m *Model) DeleteCurrentTask() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}
	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		m.Tasks = slices.Delete(m.Tasks, idx, idx+1)
	}
	newTasks := m.GetFilteredTasks()
	if m.SelectedIndex >= len(newTasks) && len(newTasks) > 0 {
		m.SelectedIndex = len(newTasks) - 1
	}
}

func (m *Model) SetDueDateForCurrentTask(dateStr string) {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}
	targetID := tasks[m.SelectedIndex].ID
	idx := m.findTaskIndexByID(targetID)
	if idx == -1 {
		return
	}
	if strings.ToLower(dateStr) == "clear" {
		m.Tasks[idx].DueDate = ""
		return
	}
	if dateStr == "" {
		return
	}
	_, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		m.ErrorMessage = "Invalid date format. Use YYYY-MM-DD"
		return
	}
	m.Tasks[idx].DueDate = dateStr
}

func (m *Model) ToggleCurrentTaskPriority() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}
	targetID := tasks[m.SelectedIndex].ID
	idx := m.findTaskIndexByID(targetID)
	if idx == -1 {
		return
	}
	priorities := []string{"", "low", "medium", "high"}
	currentPrioIdx := slices.Index(priorities, m.Tasks[idx].Priority)
	if currentPrioIdx == -1 {
		currentPrioIdx = 0
	}
	nextIdx := (currentPrioIdx + 1) % len(priorities)
	m.Tasks[idx].Priority = priorities[nextIdx]
}

func (m *Model) AddTagToCurrentTask(tag string) {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}
	targetID := tasks[m.SelectedIndex].ID
	if idx := m.findTaskIndexByID(targetID); idx != -1 {
		if !slices.Contains(m.Tasks[idx].Tags, tag) {
			m.Tasks[idx].Tags = append(m.Tasks[idx].Tags, tag)
		}
	}
}

func (m *Model) RemoveTagsFromCurrentTask() {
	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		return
	}
	targetID := tasks[m.SelectedIndex].ID
	idx := m.findTaskIndexByID(targetID)
	if idx == -1 {
		return
	}
	var newTags []string
	for j, tag := range m.Tasks[idx].Tags {
		if j < len(m.RemoveTagChecks) && !m.RemoveTagChecks[j] {
			newTags = append(newTags, tag)
		}
	}
	m.Tasks[idx].Tags = newTags
}

func (m *Model) SaveStateForUndo() {
	stateCopy := slices.Clone(m.Tasks)
	m.History = append(m.History, stateCopy)
	if len(m.History) > m.MaxHistory {
		m.History = m.History[1:]
	}
}

func (m *Model) Undo() {
	if len(m.History) == 0 {
		m.ErrorMessage = "Nothing to undo"
		return
	}
	m.Tasks = m.History[len(m.History)-1]
	m.History = m.History[:len(m.History)-1]
	m.UpdateContexts()
	m.SelectedIndex = 0
}
