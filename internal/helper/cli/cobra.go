package cli

import "github.com/spf13/cobra"

func NewCommand(use, short string, args cobra.PositionalArgs, run func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  args,
		Run:   run,
	}
}
