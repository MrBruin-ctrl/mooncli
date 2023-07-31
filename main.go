package main

import (
	"github.com/Mrbruin-ctrl/moon-cli/src/environments"
	"github.com/Mrbruin-ctrl/moon-cli/src/environments/docker"
	"github.com/c-bata/go-prompt"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	environments.Registers = map[environments.Environment]environments.Register{
		environments.Docker: docker.RegisterEnv,
	}

	environments.Logo()
	environments.Initialize()
	environments.SelectCurrentEnvironment()
	// start go-prompt as console
	p := prompt.New(
		environments.GetActiveExecutor(),
		environments.GetActiveCompleter(),
		// register hot key for select active env
		prompt.OptionAddKeyBind(environments.BuildPromptKeyBinds()...),
		// register live prefix that will be change automatically when env changed
		prompt.OptionLivePrefix(environments.LivePrefix),
		prompt.OptionTitle("MOON-CLI:一款交互式容器管理客户端"),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionOnDown(),
		prompt.OptionMaxSuggestion(8),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkBlue),
		prompt.OptionScrollbarBGColor(prompt.LightGray),
		prompt.OptionScrollbarThumbColor(prompt.Blue),
	)

	p.Run()
}
