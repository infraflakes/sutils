package find

import (
	"github.com/spf13/cobra"
)

var DirCmd = &cobra.Command{
	Use:   "dir <path> <terms...>",
	Short: "Search for directories by name",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		terms := args[1:]
		FindAndProcess(path, terms, "d", "Searching for directories with '%s' in %s\n", "Delete matched directories? (y/N): ", false)
	},
}

var DirDeleteCmd = &cobra.Command{
	Use:   "delete <path> <terms...>",
	Short: "Delete directories by name",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		terms := args[1:]
		FindAndProcess(path, terms, "d", "Searching for directories with '%s' in %s\n", "Delete matched directories? (y/N): ", true)
	},
}

func init() {
	DirCmd.AddCommand(DirDeleteCmd)
}
