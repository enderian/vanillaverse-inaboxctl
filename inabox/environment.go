package inabox

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/vanillaverse/inaboxctl/utils"
)

func Environment() map[string]string {
	contents := make(map[string]string)

	// Ports
	rootPort := viper.GetInt("deploy.root_port")
	for i := 0; i < 10; i++ {
		contents[fmt.Sprintf("INABOX_PORT_%d", i+1)] = fmt.Sprintf("%d", rootPort+i)
	}

	// Directories
	contents["INABOX_ROOT"] = viper.GetString("deploy.root")
	contents["INABOX_PROJECT"] = viper.GetString("name")
	contents["INABOX_PROJECT_DIR"] = viper.GetString("deploy.remote")
	return contents
}

func NewCompose() *utils.ComposeBuilder {
	compose := utils.NewCompose()
	compose.Env = Environment()
	return compose
}
