package docker

import (
	"docker.io/go-docker"
	"github.com/Mrbruin-ctrl/moon-cli/src/environments"
	"github.com/c-bata/go-prompt"
)

var dockerClient *docker.Client

func RegisterEnv() (*environments.Runtime, error) {
	if c, err := New(); err != nil {
		return nil, err
	} else {
		return &environments.Runtime{
			Prefix:    "docker",
			Completer: c,
			Commands:  ExtraCommands,
			Executor: environments.GetDefaultExecutor("docker", func() {
				lastQueryResult = []prompt.Suggest{}
			}),
		}, nil
	}
}
