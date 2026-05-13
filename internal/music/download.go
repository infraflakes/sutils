package music

import (
	"github.com/infraflakes/sutils/internal/helper/exec"
	"github.com/infraflakes/sutils/internal/helper/utils"
	"github.com/spf13/cobra"
)

var YTMusicDownloadCmd = &cobra.Command{
	Use:   "download [youtube-url]",
	Short: "Download audio from YouTube using yt-dlp",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		youtubeURL := args[0]
		utils.CheckErr(exec.ExecuteCommand(
			"yt-dlp",
			"--extract-audio",
			"--embed-thumbnail",
			"--add-metadata",
			"-o", "%(title)s.%(ext)s",
			youtubeURL,
		))
	},
}
