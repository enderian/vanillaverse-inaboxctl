package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanillaverse/inaboxctl/inabox"
)

// rootCmd represents the inaboxctl command
var rootCmd = &cobra.Command{
	Use:                "inaboxctl",
	Short:              "Run stuff on your inabox instance.",
	Args:               cobra.MinimumNArgs(1),
	PersistentPreRun:   before,
	Run:                proxy,
	DisableFlagParsing: true,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enables verbose output")

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(consoleCmd)
}

func before(cmd *cobra.Command, args []string) {
	// Load the global config file
	viper.SetConfigName("inabox")
	viper.SetConfigType("yml")
	viper.AddConfigPath(os.Getenv("HOME"))
	viper.AddConfigPath(os.Getenv("HOME") + "/.config")

	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	// Load the config file in the current directory (if present)
	viper.SetConfigName(".inabox")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	// Also search for the git root of the current project
	gitDirCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	if res, err := gitDirCmd.Output(); err == nil {
		viper.AddConfigPath(strings.Trim(string(res), " \n\t"))
	}

	err = viper.MergeInConfig()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("A valid .inabox.yml could not be found. Is this a VV project?")
		os.Exit(1)
	}

	projectDir := path.Dir(viper.ConfigFileUsed())
	viper.SetDefault("name", path.Base(projectDir))

	// Set the default config values
	viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose"))

	viper.SetDefault("deploy.root", fmt.Sprintf("/srv/inabox/%s", os.Getenv("USER")))
	viper.SetDefault("deploy.local", projectDir)
	viper.SetDefault("deploy.remote", path.Join(viper.GetString("deploy.root"), "projects", viper.GetString("name")))
	viper.SetDefault("deploy.port", 22)
	viper.SetDefault("deploy.exclude_files", []string{".git", ".gradle", ".bundle", "build", "node_modules"})

	viper.SetDefault("compose.file", "docker-compose.yml")
	viper.SetDefault("compose.services", "docker-compose.services.yml")
	viper.SetDefault("compose.project", fmt.Sprintf("inabox_%s", os.Getenv("USER")))
}

func proxy(cmd *cobra.Command, args []string) {
	err := inabox.Proxy(args)
	cobra.CheckErr(err)
}
