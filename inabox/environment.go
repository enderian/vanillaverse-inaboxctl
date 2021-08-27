package inabox

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/vanillaverse/inaboxctl/utils"
)

var redisPort = "6379"

func Environment() map[string]string {
	contents := make(map[string]string)

	// Ports
	rootPort := viper.GetInt("deploy.root_port")
	for _, service := range viper.GetStringSlice("deploy.services") {
		contents[fmt.Sprintf("INABOX_%s_PORT", strings.ToUpper(service))] = fmt.Sprintf("%d", rootPort)
		rootPort += 1
	}
	for i := 1; i <= 9; i++ {
		contents[fmt.Sprintf("INABOX_PORT_%d", i)] = fmt.Sprintf("%d", rootPort)
		rootPort += 1
	}

	contents["INABOX_HOST"] = viper.GetString("deploy.host")
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
