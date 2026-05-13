package cmd

import (
	"github.com/infraflakes/sutils/cmd/cd"
	"github.com/infraflakes/sutils/cmd/find"
	"github.com/infraflakes/sutils/cmd/music"
	"github.com/infraflakes/sutils/cmd/todo"
	"github.com/infraflakes/sutils/cmd/zip"
	"github.com/infraflakes/sutils/internal/helper/exec"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sn",
	Short: "Sane Utils — a suite of CLI utilities",
	Long:  `Sane Utils is an opinionated CLI suite to streamline many command line work.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	_ = rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(cd.RootCmd)
	rootCmd.AddCommand(find.RootCmd)
	rootCmd.AddCommand(music.RootCmd)
	rootCmd.AddCommand(todo.RootCmd)
	rootCmd.AddCommand(zip.RootCmd)
	rootCmd.PersistentFlags().BoolVar(&exec.DryRun, "dry-run", false, "print the command that would be executed instead of executing it")
}
