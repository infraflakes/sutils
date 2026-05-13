package find

import (
	"github.com/spf13/cobra"
)

var FileCmd = &cobra.Command{
	Use:   "file <terms...> [path]",
	Short: "Search for files by name",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		terms, path := parseTermsPath(args)
		FindAndProcess(path, terms, "f", "Searching for files with '%s' in %s\n", "Delete matched files? (y/N): ", false)
	},
}

var FileDeleteCmd = &cobra.Command{
	Use:   "delete <terms...> [path]",
	Short: "Delete files by name",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		terms, path := parseTermsPath(args)
		FindAndProcess(path, terms, "f", "Searching for files with '%s' in %s\n", "Delete matched files? (y/N): ", true)
	},
}

func init() {
	FileCmd.AddCommand(FileDeleteCmd)
}
