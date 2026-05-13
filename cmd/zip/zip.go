package zip

import (
	zip "github.com/infraflakes/sutils/internal/zip"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zip",
	Short: "Archive and extract files with 7z",
}

func init() {
	RootCmd.AddCommand(zip.ZipCmd)
	RootCmd.AddCommand(zip.UnzipCmd)
}
