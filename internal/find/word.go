package find

import (
	"fmt"

	"github.com/infraflakes/sutils/internal/helper/cli"
	"github.com/infraflakes/sutils/internal/helper/exec"
	"github.com/spf13/cobra"
)

var WordCmd = &cobra.Command{
	Use:   "word <terms...> [path]",
	Short: "Search for words inside files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		terms, path := parseTermsPath(args)

		for _, term := range terms {
			fmt.Printf("Searching for '%s' in %s\n", term, path)
			output, err := exec.RunCommand("grep", "-rE", term, path)
			if err != nil {
				fmt.Printf("Error searching for '%s': %v\n", term, err)
				continue
			}
			for _, line := range output {
				fmt.Println(line)
			}
		}
	},
}

var WordDeleteCmd = &cobra.Command{
	Use:   "delete <terms...> [path]",
	Short: "Delete files containing matching words",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		terms, path := parseTermsPath(args)

		for _, term := range terms {
			fmt.Printf("Searching for '%s' in %s\n", term, path)
			output, err := exec.RunCommand("grep", "-rE", term, path)
			if err != nil {
				fmt.Printf("Error searching for '%s': %v\n", term, err)
				continue
			}
			for _, line := range output {
				fmt.Println(line)
			}

			if cli.Confirm("Delete matched files? (y/N): ") {
				DeleteGrepMatches(path, term)
			}
		}
	},
}

func init() {
	WordCmd.AddCommand(WordDeleteCmd)
}
