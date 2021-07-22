package inabox

import (
	"fmt"

	"github.com/vanillaverse/inaboxctl/vssh"
)

// Proxy proxies the command given to the CLI to the remote docker-compose project.
func Proxy(args []string) error {
	ssh, err := vssh.NewSession()
	if err != nil {
		return fmt.Errorf("error while creating ssh session: %v", err)
	}

	compose := NewCompose()
	compose.WithProject()
	compose.WithServices()
	return ssh.Run(compose.Command(args...))
}
