package zip

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var ZipCmd = &cobra.Command{
	Use:   "zip [archive-name] [target-to-archive]",
	Short: "Archive files with 7z",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		archiveName := args[0]
		targets := ExpandTargets(args[1:])
		BuildArchiveCommand(archiveName, targets, "")
	},
}

var zipPasswordCmd = &cobra.Command{
	Use:   "password [archive-name] [target-to-archive]",
	Short: "Archive files with 7z and a password",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		archiveName := args[0]
		targets := ExpandTargets(args[1:])
		fmt.Print("Enter archive password: ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("Failed to read password:", err)
			return
		}
		password := string(bytePassword)
		if password == "" {
			fmt.Println("Password cannot be empty.")
			return
		}
		BuildArchiveCommand(archiveName, targets, password)
	},
}

func init() {
	ZipCmd.AddCommand(zipPasswordCmd)
}
