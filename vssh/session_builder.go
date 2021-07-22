package vssh

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func NewSession() (*Session, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	agentConn, err := net.Dial("unix", socket)

	if err != nil {
		fmt.Printf("Failed to open SSH_AUTH_SOCK: %v", err)
		os.Exit(1)
	}

	agentClient := agent.NewClient(agentConn)
	config := &ssh.ClientConfig{
		User:            viper.GetString("deploy.user"),
		Auth:            []ssh.AuthMethod{ssh.PublicKeysCallback(agentClient.Signers)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	config.SetDefaults()
	addr := fmt.Sprintf("%s:%d", viper.GetString("deploy.host"), viper.GetInt("deploy.port"))
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("error while establishing ssh connection: %v", err)
	}

	return &Session{conn}, err
}
