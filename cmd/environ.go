package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vanillaverse/inaboxctl/inabox"
)

var environCmd = &cobra.Command{
	Use:   "environ",
	Short: "Print out the enrivonment variable",
	Run:   environRun,
}

func environRun(cmd *cobra.Command, args []string) {
	env := inabox.Environment()
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("%s=%s\n", key, env[key])
	}
}
