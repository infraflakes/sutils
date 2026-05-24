package cd

import (
	"fmt"
	"os"

	cdinit "github.com/infraflakes/sutils/internal/cd"
	cdalias "github.com/infraflakes/sutils/internal/cd/alias"
	cdtui "github.com/infraflakes/sutils/internal/cd/tui"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cd [alias]",
	Short: "Switch directories using aliases and TUI",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var target string
		if len(args) == 0 {
			selected, err := cdtui.RunTUI()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if selected == "" {
				os.Exit(0)
			}
			target = selected
		} else {
			target = args[0]
		}

		path, err := cdalias.Priority(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: '%s' is not a valid path or alias\n", target)
			os.Exit(1)
		}

		fmt.Print(path)
	},
}

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage directory aliases",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var aliasAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add current directory as an alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := cdalias.AddAlias(name, cwd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Added alias: %s -> %s\n", name, cwd)
	},
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		aliases, err := cdalias.ReadAliases()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(aliases) == 0 {
			fmt.Fprintln(os.Stderr, "No aliases configured.")
			return
		}

		for k, v := range aliases {
			fmt.Printf("%s = %s\n", k, v)
		}
	},
}

var aliasExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export aliases to the current directory",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		dest := cwd + "/" + cdalias.ConfigFileName
		if err := cdalias.ExportAliases(dest); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Exported aliases to: %s\n", dest)
	},
}

var aliasDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := cdalias.RemoveAlias(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Deleted alias: %s\n", name)
	},
}

var aliasWipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Wipe all aliases",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cdalias.WipeAliases(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "All aliases wiped.")
	},
}

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell initialization script (fish, bash, zsh)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		script, err := cdinit.GenerateInit(shell)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(script)
	},
}

func init() {
	aliasCmd.AddCommand(aliasAddCmd)
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasDeleteCmd)
	aliasCmd.AddCommand(aliasWipeCmd)
	aliasCmd.AddCommand(aliasExportCmd)
	RootCmd.AddCommand(aliasCmd)
	RootCmd.AddCommand(initCmd)
}
