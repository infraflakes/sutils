package zip

import (
	"github.com/spf13/cobra"
)

var UnzipCmd = &cobra.Command{
	Use:   "unzip [target-to-unarchive]",
	Short: "Extract archives with 7z",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		BuildExtractCommand(args[0], "")
	},
}

var unzipPasswordCmd = &cobra.Command{
	Use:   "password [password] [target-to-unarchive]",
	Short: "Extract archives with a password",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		BuildExtractCommand(args[1], args[0])
	},
}

func init() {
	UnzipCmd.AddCommand(unzipPasswordCmd)
}
