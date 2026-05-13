package find

import (
	"github.com/spf13/cobra"
)

var DirCmd = &cobra.Command{
	Use:   "dir <terms...> [path]",
	Short: "Search for directories by name",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		terms, path := parseTermsPath(args)
		FindAndProcess(path, terms, "d", "Searching for directories with '%s' in %s\n", "Delete matched directories? (y/N): ", false)
	},
}

var DirDeleteCmd = &cobra.Command{
	Use:   "delete <terms...> [path]",
	Short: "Delete directories by name",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		terms, path := parseTermsPath(args)
		FindAndProcess(path, terms, "d", "Searching for directories with '%s' in %s\n", "Delete matched directories? (y/N): ", true)
	},
}

func init() {
	DirCmd.AddCommand(DirDeleteCmd)
}
