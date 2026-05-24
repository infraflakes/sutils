package zip

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
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
	Use:   "password [target-to-unarchive]",
	Short: "Extract archives with a password",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
		BuildExtractCommand(args[0], password)
	},
}

func init() {
	UnzipCmd.AddCommand(unzipPasswordCmd)
}
