package find

import (
	find "github.com/infraflakes/sutils/internal/find"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "find",
	Short: "Search for words, files, or directories",
}

func init() {
	RootCmd.AddCommand(find.WordCmd)
	RootCmd.AddCommand(find.FileCmd)
	RootCmd.AddCommand(find.DirCmd)
}
