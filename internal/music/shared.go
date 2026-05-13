package music

import "strings"

func FormatPlaylistLines(lines []string) []string {
	var formatted []string
	for _, line := range lines {
		line = strings.ReplaceAll(line, "/", `\`) + "\r"
		formatted = append(formatted, line)
	}
	return formatted
}
