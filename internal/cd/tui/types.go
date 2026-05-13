package tui

type entry struct {
	name  string
	isDir bool
	path  string
}

type model struct {
	currentDir     string
	selectedIdx    int
	parentEntries  []entry
	currentEntries []entry
	previewEntries []entry
	width, height  int
	quitting       bool
	finalPath      string
	showFiles      bool
	showHidden     bool
	searchMode     bool
	searchQuery    string
}
