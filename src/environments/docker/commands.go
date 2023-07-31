package docker

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Mrbruin-ctrl/moon-cli/src/environments"
	"github.com/Mrbruin-ctrl/moon-cli/src/internal/util"
	"github.com/Mrbruin-ctrl/moon-cli/src/suggestions"
	"io"
	"strings"
)

var ExtraCommands = []environments.Command{
	{
		Text:         "x-batch",
		Provider:     DockerType,
		Description:  "批量管理容器和镜像",
		Environments: []environments.Environment{environments.Docker},
		MatchFn:      environments.StartWithMatch,
		Fn: func(args []string, writer io.Writer) {
			if len(args) < 2 || (args[1] != "container" && args[1] != "image") {
				_, _ = writer.Write([]byte("需要提供 'container' or 'image' 类型参数 \n" +
					"请使用 \"x-batch container\" 或者 \"x-batch image\" 进行重试\n"))
				return
			}

			resources := map[string]struct {
				fn  suggestions.Provider
				ops []string
			}{
				"container": {
					fn:  provideContainerSuggestion,
					ops: []string{"rm", "rm -f", "start", "stop", "restart", "pause", "unpause", "kill"},
				},
				"image": {
					fn:  provideImagesSuggestion,
					ops: []string{"rmi", "rmi -f", "create", "run"},
				},
			}

			resource := resources[args[1]]
			var opts []string
			col := int(util.GetWindowWidth() - 15)
			for _, v := range resource.fn() {
				opts = append(opts, v.Text+" | "+util.SubString(v.Description, 0, col-len(v.Text)))
			}
			var qs = []*survey.Question{
				{
					Name: "objective",
					Prompt: &survey.MultiSelect{
						Message:  "操作哪些资源 ?",
						Options:  opts,
						PageSize: 25,
					},
					Validate: func(val interface{}) error {
						if ans, ok := val.([]survey.OptionAnswer); !ok || len(ans) == 0 {
							return errors.New("请选择至少一种资源")
						}
						return nil
					},
				},
				{
					Name: "operation",
					Prompt: &survey.Select{
						Message: "想做什么操作 ?",
						Options: resource.ops,
					},
				},
				{
					Name: "confirm",
					Prompt: &survey.Confirm{
						Message: "确认继续 ?",
					},
				},
			}
			answers := struct {
				Objective []string
				Operation string
				Confirm   bool
			}{}

			err := survey.Ask(qs, &answers, survey.WithKeepFilter(true))
			if err == terminal.InterruptErr {
				_, _ = writer.Write([]byte("operation interrupted\n"))
			} else if err != nil {
				_, _ = writer.Write([]byte(err.Error() + "\n"))
			}

			if answers.Confirm {
				for _, o := range answers.Objective {
					environments.Executor("docker", answers.Operation+" "+strings.TrimSpace(strings.Split(o, "|")[0]))
				}
			}
		},
	},
}
