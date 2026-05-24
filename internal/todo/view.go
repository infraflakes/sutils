package todo

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	renderer = lipgloss.NewRenderer(os.Stderr)

	baseStyle = renderer.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	titleStyle = renderer.NewStyle().
			Bold(true).
			Reverse(true).
			Padding(0, 1)

	taskStyle = renderer.NewStyle().
			PaddingLeft(2)

	selectedTaskStyle = renderer.NewStyle().
				Bold(true).
				Reverse(true).
				PaddingLeft(2)

	completedTaskStyle = renderer.NewStyle().
				Faint(true).
				Strikethrough(true)

	highPriorityStyle = renderer.NewStyle().
				Bold(true)

	mediumPriorityStyle = renderer.NewStyle()

	lowPriorityStyle = renderer.NewStyle().
				Faint(true)

	contextStyle = renderer.NewStyle().
			Bold(true).
			Underline(true)

	errorStyle = renderer.NewStyle().
			Bold(true).
			Reverse(true)

	inputStyle = renderer.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1).
			Margin(1)
)

func (m Model) View() string {
	if m.HelpVisible {
		return m.renderFullHelpView()
	}

	switch m.ViewMode {
	case InputView:
		return m.RenderInputView()
	case DateInputView:
		return m.RenderDateInputView()
	case RemoveTagView:
		return m.RenderRemoveTagView()
	case KanbanView:
		return m.RenderKanbanView()
	case StatsView:
		return m.RenderStatsView()
	default:
		return m.RenderNormalView()
	}
}

func (m Model) renderFullHelpView() string {
	helpBoxStyle := renderer.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	m.Help.ShowAll = true
	helpContent := m.Help.View(m.KeyMap)
	titledHelp := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Keybindings"),
		helpContent,
	)

	return lipgloss.Place(m.WindowWidth, m.WindowHeight, lipgloss.Center, lipgloss.Center, helpBoxStyle.Render(titledHelp))
}

func (m Model) RenderNormalView() string {
	var mainContent strings.Builder

	contextText := fmt.Sprintf("Context: %s", m.CurrentContext)
	mainContent.WriteString(titleStyle.Render(contextText) + "\n\n")

	tasks := m.GetFilteredTasks()
	if len(tasks) == 0 {
		if len(m.Contexts) == 0 {
			mainContent.WriteString("No contexts exist. Press 'n' to create one.\n")
		} else {
			mainContent.WriteString("No tasks in this context. Press 'a' to add one.\n")
		}
	} else {
		for i, task := range tasks {
			taskLine := m.RenderTask(task, i == m.SelectedIndex, m.MovingMode && task.ID == m.MovingTaskID)
			mainContent.WriteString(taskLine + "\n")
		}
	}

	if m.ErrorMessage != "" {
		mainContent.WriteString("\n" + errorStyle.Render(m.ErrorMessage) + "\n")
	}

	return baseStyle.Render(mainContent.String())
}

func (m Model) RenderTask(task Task, selected, moving bool) string {
	checkbox := "[ ]"
	if task.Checked {
		checkbox = "[✓]"
	}

	priority := ""
	switch task.Priority {
	case "high":
		priority = highPriorityStyle.Render("!!! ")
	case "medium":
		priority = mediumPriorityStyle.Render("!! ")
	case "low":
		priority = lowPriorityStyle.Render("! ")
	}

	taskText := task.Task

	tags := ""
	if len(task.Tags) > 0 {
		tags = " > " + strings.Join(task.Tags, ", ")
	}

	dueDate := ""
	if task.DueDate != "" {
		dueDate = fmt.Sprintf(" [Due: %s]", task.DueDate)
	}

	text := fmt.Sprintf("%s %s%s%s", checkbox, taskText, tags, dueDate)

	style := taskStyle
	if task.Checked {
		style = completedTaskStyle
	}

	if selected {
		style = style.Reverse(true).Bold(true)
	}

	if moving {
		style = style.Bold(true)
	}

	return priority + style.Render(text)
}

func (m Model) RenderInputView() string {
	content := inputStyle.Render(
		fmt.Sprintf("%s\n\n%s", m.InputPrompt, m.TextInput.View()),
	)
	return lipgloss.Place(m.WindowWidth, m.WindowHeight, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) RenderDateInputView() string {
	var content strings.Builder
	content.WriteString("Set due date (YYYY-MM-DD):\n\n")
	inputs := []string{
		fmt.Sprintf("Day: %s", m.DateInputs[0].View()),
		fmt.Sprintf("Month: %s", m.DateInputs[1].View()),
		fmt.Sprintf("Year: %s", m.DateInputs[2].View()),
	}
	for i, input := range inputs {
		if i == m.DateInputIndex {
			content.WriteString(selectedTaskStyle.Render(input) + "\n")
		} else {
			content.WriteString(input + "\n")
		}
	}
	return inputStyle.Render(content.String())
}

func (m Model) RenderRemoveTagView() string {
	var content strings.Builder
	content.WriteString("Select tags to remove:\n\n")
	task := m.GetCurrentTask()
	for i, tag := range task.Tags {
		checkbox := "[ ]"
		if m.RemoveTagChecks[i] {
			checkbox = "[✓]"
		}
		line := fmt.Sprintf("%s %s", checkbox, tag)
		if i == m.RemoveTagIndex {
			content.WriteString(selectedTaskStyle.Render(line) + "\n")
		} else {
			content.WriteString(line + "\n")
		}
	}
	return inputStyle.Render(content.String())
}

func (m Model) RenderKanbanView() string {
	var content strings.Builder
	title := titleStyle.Render("Kanban View (←/→/↑/↓ scroll, esc to return)")
	content.WriteString(title + "\n")

	if len(m.Contexts) == 0 {
		content.WriteString("No contexts available.\n")
		return baseStyle.Render(content.String())
	}

	const (
		fixedColWidth  = 35
		separatorWidth = 3
	)

	numVisibleCols := max(m.WindowWidth/(fixedColWidth+separatorWidth), 1)

	if m.KanbanScrollX > len(m.Contexts)-numVisibleCols {
		m.KanbanScrollX = max(0, len(m.Contexts)-numVisibleCols)
	}
	if m.KanbanScrollX < 0 {
		m.KanbanScrollX = 0
	}

	startCol := m.KanbanScrollX
	endCol := min(startCol+numVisibleCols, len(m.Contexts))
	visibleContexts := m.Contexts[startCol:endCol]

	columnStyle := lipgloss.NewStyle().Width(fixedColWidth).Padding(0, 1)
	taskTextStyle := lipgloss.NewStyle().Width(fixedColWidth - 2)

	var columns []string
	for _, context := range visibleContexts {
		var column strings.Builder
		header := contextStyle.Render(context)
		column.WriteString(header + "\n")
		column.WriteString(strings.Repeat("─", fixedColWidth) + "\n")

		tasks := m.GetTasksForContext(context)
		for _, task := range tasks {
			var taskLine strings.Builder
			if task.Checked {
				taskLine.WriteString("✓ ")
			} else {
				taskLine.WriteString("• ")
			}
			fullTaskText := task.Task
			if len(task.Tags) > 0 {
				fullTaskText += " > " + strings.Join(task.Tags, ", ")
			}
			if task.DueDate != "" {
				fullTaskText += fmt.Sprintf(" [Due: %s]", task.DueDate)
			}
			wrappedText := taskTextStyle.Render(fullTaskText)
			if task.Checked {
				taskLine.WriteString(completedTaskStyle.Render(wrappedText))
			} else {
				taskLine.WriteString(wrappedText)
			}
			column.WriteString(taskLine.String() + "\n")
		}
		columns = append(columns, columnStyle.Render(column.String()))
	}

	board := lipgloss.JoinHorizontal(lipgloss.Top, columns...)
	boardLines := strings.Split(board, "\n")

	top := m.KanbanScrollY
	bottom := top + m.WindowHeight - lipgloss.Height(title) - 1
	if top < 0 {
		top = 0
	}
	if bottom > len(boardLines) {
		bottom = len(boardLines)
	}
	if top > bottom {
		top = max(0, bottom-m.WindowHeight)
		m.KanbanScrollY = top
	}

	visibleLines := boardLines[top:bottom]
	content.WriteString(strings.Join(visibleLines, "\n"))

	return baseStyle.Render(content.String())
}

func (m Model) RenderStatsView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("Statistics (ESC to return)") + "\n\n")

	total := len(m.Tasks)
	completed := 0
	for _, task := range m.Tasks {
		if task.Checked {
			completed++
		}
	}

	completionRate := 0.0
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	fmt.Fprintf(&content, "Total Tasks: %d\n", total)
	fmt.Fprintf(&content, "Completed: %d (%.1f%%)\n\n", completed, completionRate)

	content.WriteString("Context Statistics:\n")
	for _, context := range m.Contexts {
		tasks := m.GetTasksForContext(context)
		ctxTotal := len(tasks)
		ctxCompleted := 0
		for _, task := range tasks {
			if task.Checked {
				ctxCompleted++
			}
		}

		ctxRate := 0.0
		if ctxTotal > 0 {
			ctxRate = float64(ctxCompleted) / float64(ctxTotal) * 100
		}

		fmt.Fprintf(&content, "  %s: %d/%d (%.1f%%)\n",
			contextStyle.Render(context), ctxCompleted, ctxTotal, ctxRate)
	}

	return baseStyle.Render(content.String())
}
