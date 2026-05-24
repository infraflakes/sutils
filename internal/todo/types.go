package todo

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
)

type Task struct {
	ID       int      `json:"id"`
	Task     string   `json:"task"`
	Checked  bool     `json:"checked"`
	Context  string   `json:"context"`
	Priority string   `json:"priority,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	DueDate  string   `json:"due_date,omitempty"`
}

type ViewMode int

const (
	NormalView ViewMode = iota
	KanbanView
	StatsView
	InputView
	DateInputView
	RemoveTagView
)

type InputMode int

const (
	AddTaskInput InputMode = iota
	EditTaskInput
	AddContextInput
	RenameContextInput
	AddTagInput
	DeleteConfirmInput
)

type Model struct {
	Tasks          []Task
	Contexts       []string
	CurrentContext string
	SelectedIndex  int
	NextID         int

	ViewMode  ViewMode
	InputMode InputMode

	MovingMode    bool
	MovingTaskID  int
	KanbanScrollY int
	KanbanScrollX int

	TextInput       textinput.Model
	DateInputs      []textinput.Model
	DateInputIndex  int
	RemoveTagIndex  int
	RemoveTagChecks []bool
	InputPrompt     string

	WindowWidth  int
	WindowHeight int
	ErrorMessage string

	History    [][]Task
	MaxHistory int

	KeyMap      KeyMap
	Help        help.Model
	HelpVisible bool

	ConfigFilePath string
}

type KeyMap struct {
	Up             key.Binding
	Down           key.Binding
	Left           key.Binding
	Right          key.Binding
	Toggle         key.Binding
	Add            key.Binding
	Edit           key.Binding
	Delete         key.Binding
	AddContext     key.Binding
	RenameContext  key.Binding
	DeleteContext  key.Binding
	TogglePriority key.Binding
	AddTag         key.Binding
	RemoveTag      key.Binding
	SetDueDate     key.Binding
	ClearDueDate   key.Binding
	KanbanView     key.Binding
	StatsView      key.Binding
	Undo           key.Binding
	Move           key.Binding
	Help           key.Binding
	Quit           key.Binding
	Back           key.Binding
	Enter          key.Binding
	Nav            key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "prev context"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next context"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		AddContext: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new context"),
		),
		RenameContext: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename context"),
		),
		DeleteContext: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "delete context"),
		),
		TogglePriority: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "priority"),
		),
		AddTag: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "add tag"),
		),
		RemoveTag: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "remove tag"),
		),
		SetDueDate: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "due date"),
		),
		ClearDueDate: key.NewBinding(
			key.WithKeys("U"),
			key.WithHelp("U", "clear due"),
		),
		KanbanView: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "kanban"),
		),
		StatsView: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stats"),
		),
		Undo: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp("z", "undo"),
		),
		Move: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "move"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Nav: key.NewBinding(
			key.WithKeys("↑", "↓", "←", "→"),
			key.WithHelp("↑↓←→", "navigation"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Nav, k.Toggle, k.Add, k.Edit, k.Delete, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav},
		{k.Toggle, k.Add, k.Edit, k.Delete, k.Move},
		{k.AddContext, k.RenameContext, k.DeleteContext},
		{k.TogglePriority, k.AddTag, k.RemoveTag, k.SetDueDate, k.ClearDueDate},
		{k.KanbanView, k.StatsView},
		{k.Undo, k.Help, k.Back, k.Quit},
	}
}
