package inabox

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/vanillaverse/inaboxctl/vssh"
)

// Console creates a console for the specified service
func Console(service string) error {
	ssh, err := vssh.NewSession()
	if err != nil {
		return fmt.Errorf("error while creating ssh session: %v", err)
	}

	compose := NewCompose()
	compose.WithProject()
	container, err := ssh.Output(compose.Command("ps", "-q", service))
	if err != nil {
		return fmt.Errorf("error while getting container: %v", err)
	}

	containerClean := strings.Trim(string(container), " \n\t")
	if viper.GetBool("verbose") {
		fmt.Printf("will attach to container: %s\n", containerClean)
	}

	fmt.Printf("Attaching to %s...\n", service)
	return ssh.Console("docker attach " + containerClean)
}
