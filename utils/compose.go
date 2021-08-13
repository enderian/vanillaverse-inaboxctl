package utils

import (
	"path"
	"strings"

	"github.com/spf13/viper"
)

type ComposeBuilder struct {
	projectName         string
	servicesComposeFile string
	composeFile         string

	Services []string
	Env      map[string]string
}

func NewCompose() *ComposeBuilder {
	return &ComposeBuilder{
		projectName: viper.GetString("compose.project"),
	}
}

func (b *ComposeBuilder) WithServices() {
	b.servicesComposeFile = path.Join(viper.GetString("deploy.root"), viper.GetString("compose.services"))
}

func (b *ComposeBuilder) WithProject() {
	b.composeFile = path.Join(viper.GetString("deploy.remote"), viper.GetString("compose.file"))
}

func (b *ComposeBuilder) Command(arguments ...string) string {
	var parts []string
	for k, v := range b.Env {
		parts = append(parts, k+"="+v)
	}
	parts = append(parts, "docker-compose", "-p", b.projectName)
	if b.servicesComposeFile != "" {
		parts = append(parts, "-f", b.servicesComposeFile)
	}
	if b.composeFile != "" {
		parts = append(parts, "-f", b.composeFile)
	}
	parts = append(parts, arguments...)
	parts = append(parts, b.Services...)
	return strings.Join(parts, " ")
}

func ParseServices(services []byte) []string {
	var servicesList []string
	for _, service := range strings.Split(string(services), "\n") {
		if service != "" {
			servicesList = append(servicesList, strings.TrimSpace(service))
		}
	}
	return servicesList
}
