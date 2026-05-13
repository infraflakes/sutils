package find

import (
	"github.com/spf13/cobra"
)

var FileCmd = &cobra.Command{
	Use:   "file <path> <terms...>",
	Short: "Search for files by name",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		terms := args[1:]
		FindAndProcess(path, terms, "f", "Searching for files with '%s' in %s\n", "Delete matched files? (y/N): ", false)
	},
}

var FileDeleteCmd = &cobra.Command{
	Use:   "delete <path> <terms...>",
	Short: "Delete files by name",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		terms := args[1:]
		FindAndProcess(path, terms, "f", "Searching for files with '%s' in %s\n", "Delete matched files? (y/N): ", true)
	},
}

func init() {
	FileCmd.AddCommand(FileDeleteCmd)
}
