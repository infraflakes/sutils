package music

import (
	music "github.com/infraflakes/sutils/internal/music"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "music",
	Short: "Music utilities for conversion and download",
}

func init() {
	RootCmd.AddCommand(music.ConvertCmd)
	RootCmd.AddCommand(music.YTMusicDownloadCmd)
}
