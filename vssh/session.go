package vssh

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type Session struct {
	conn *ssh.Client
}

// Console creates a pseudoconsole by running the provided command
func (s *Session) Console(command string) error {
	session, err := s.conn.NewSession()
	defer session.Close()

	if err != nil {
		return fmt.Errorf("could not create new session: %v", err)
	}
	if viper.GetBool("verbose") {
		fmt.Printf("executing (ssh,console): %s\n", command)
	}

	// Get a pseudo-terminal
	if err := session.RequestPty(fmt.Sprintf("xterm-256color"), 80, 40, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %v", err)
	}

	// Setup all the pipes!
	pipeIn, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("could not get stdin pipe: %v", err)
	}
	go io.Copy(pipeIn, os.Stdin)
	pipeOut, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not get stdout pipe: %v", err)
	}
	go io.Copy(os.Stdout, pipeOut)
	pipeErr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not get stderr pipe: %v", err)
	}
	go io.Copy(os.Stderr, pipeErr)

	// Execute the command
	return session.Run(command)
}

func (s *Session) Output(command string) ([]byte, error) {
	session, err := s.conn.NewSession()
	defer session.Close()

	if err != nil {
		return nil, fmt.Errorf("could not create new session: %v", err)
	}
	if viper.GetBool("verbose") {
		fmt.Printf("executing (ssh,output): %s\n", command)
	}
	return session.Output(command)
}

func (s *Session) Run(command string) error {
	session, err := s.conn.NewSession()
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err != nil {
		return fmt.Errorf("could not create new session: %v", err)
	}
	if viper.GetBool("verbose") {
		fmt.Printf("executing (ssh): %s\n", command)
	}
	return session.Run(command)
}
