package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanillaverse/inaboxctl/inabox"
)

var buildFlag bool
var pullFlag bool
var detachFlag bool
var noHooksFlag bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your changes",
	Run:   deployRun,
}

func init() {
	deployCmd.PersistentFlags().BoolVarP(&buildFlag, "build", "b", false, "Unconditionally build the project")
	deployCmd.PersistentFlags().BoolVarP(&pullFlag, "pull", "p", false, "Unconditionally pull the project's images")
	deployCmd.PersistentFlags().BoolVarP(&detachFlag, "detach", "d", false, "Run in the background")
	deployCmd.PersistentFlags().BoolVar(&noHooksFlag, "no-hooks", false, "Skip hooks for this deployment")
}

func deployRun(cmd *cobra.Command, args []string) {
	err := inabox.Deploy(args, pullFlag, buildFlag, detachFlag, noHooksFlag)
	cobra.CheckErr(err)
}
