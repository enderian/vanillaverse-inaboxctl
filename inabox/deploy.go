package inabox

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
	"github.com/vanillaverse/inaboxctl/utils"
	"github.com/vanillaverse/inaboxctl/vssh"
)

// Deploy runs the deploy pipeline to inabox.
func Deploy(services []string, pull bool, build bool, detach bool, no_hooks bool) error {
	fmt.Println("Syncing files with inabox...")
	if err := syncFiles(); err != nil {
		return fmt.Errorf("error while syncing files: %v", err)
	}

	ssh, err := vssh.NewSession()
	if err != nil {
		return fmt.Errorf("error while creating ssh session: %v", err)
	}

	fmt.Println("Making sure services are up...")
	err = startServices(ssh)
	if err != nil {
		return fmt.Errorf("error while starting services: %v", err)
	}

	compose := NewCompose()
	compose.WithProject()

	// Configure the containers to (build,) start and attach to.
	configServices := viper.GetStringSlice("project.default_services")
	if len(services) > 0 {
		compose.Services = services
	} else if configServices != nil {
		compose.Services = configServices
	} else {
		srv, err := ssh.Output(compose.Command("ps", "--services"))
		if err == nil {
			compose.Services = utils.ParseServices(srv)
		}
	}

	// Error out if no services
	if compose.Services == nil || len(compose.Services) < 1 {
		return fmt.Errorf("no service containers detected, please specify them explicitly")
	}

	// Stop services first
	fmt.Printf("Stopping %s...\n", strings.Join(compose.Services, ", "))
	err = ssh.Run(compose.Command("stop"))
	if err != nil {
		return fmt.Errorf("error while stopping services: %v", err)
	}

	// Run pre-deploy hooks
	if hook := viper.GetString("deploy.pre_hook"); !no_hooks && hook != "" {
		fmt.Println("Running pre-deploy hook...")
		if err := ssh.Run(fmt.Sprintf("cd %s && %s", viper.GetString("deploy.remote"), hook)); err != nil {
			return fmt.Errorf("error while running pre-deploy hook: %v", err)
		}
	}

	// Pull containers conditionally.
	if pull {
		fmt.Println("Pulling images and dependencies...")
		err = ssh.Run(compose.Command("pull", "-q", "--ignore-pull-failures", "--include-deps"))
		if err != nil {
			return fmt.Errorf("error while building services: %v", err)
		}
	}

	// Build the containers if necessary.
	if build {
		fmt.Printf("\rBuilding containers...")
		err = ssh.Run(compose.Command("build"))
		if err != nil {
			return fmt.Errorf("error while building services: %v", err)
		}
	}

	fmt.Printf("Deploying %s...\n", strings.Join(compose.Services, ", "))
	compose.WithServices()
	err = ssh.Run(compose.Command("up", "-d", "--force-recreate", "--remove-orphans"))
	if err != nil {
		return fmt.Errorf("error while deploying services: %v", err)
	}

	// Attach to containers for logging.
	if !detach {
		err = ssh.Run(compose.Command("logs", "-f"))
		if err != nil {
			return fmt.Errorf("error while requesting logs: %v", err)
		}
	}
	return nil
}

// syncFiles syncs the files between the local project directory and inabox.
func syncFiles() error {
	args := []string{"--delete", "--progress", "-arzh"}
	for _, excluded := range viper.GetStringSlice("deploy.exclude_files") {
		args = append(args, "--exclude", excluded)
	}
	args = append(args, fmt.Sprintf(
		"%s/",
		viper.GetString("deploy.local"),
	))
	args = append(args, fmt.Sprintf(
		"%s@%s:%s",
		viper.GetString("deploy.user"),
		viper.GetString("deploy.host"),
		viper.GetString("deploy.remote"),
	))

	rsync := exec.Command("rsync", args...)

	// Verbose logging
	if viper.GetBool("verbose") {
		fmt.Printf("executing: %s\n", strings.Join(rsync.Args, " "))
	}

	rsync.Stdout = os.Stdout
	rsync.Stderr = os.Stderr
	if err := rsync.Run(); err != nil {
		return fmt.Errorf("error while syncing files: %v", err)
	}
	return nil
}

// startServices starts any stopped service containers.
func startServices(session *vssh.Session) error {
	compose := NewCompose()
	compose.WithServices()
	compose.Env["COMPOSE_IGNORE_ORPHANS"] = "1"
	return session.Run(compose.Command("up", "-d", "--no-recreate"))
}
