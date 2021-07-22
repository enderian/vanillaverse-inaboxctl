package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanillaverse/inaboxctl/inabox"
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Attach to the console of the given service",
	Args:  cobra.ExactArgs(1),
	Run:   consoleRun,
}

func consoleRun(cmd *cobra.Command, args []string) {
	err := inabox.Console(args[0])
	cobra.CheckErr(err)
}
