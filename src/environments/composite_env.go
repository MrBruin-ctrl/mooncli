package environments

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mrbruin-ctrl/moon-cli/lib/go-ansi"
	"github.com/Mrbruin-ctrl/moon-cli/src/internal/template"
	"github.com/Mrbruin-ctrl/moon-cli/src/suggestions"
	prompt "github.com/c-bata/go-prompt"
	"os"
	"time"
)

type Runtime struct {
	Prefix     string
	Completer  *suggestions.GenericCompleter
	Executor   prompt.Executor
	Setup      func()
	LivePrefix func() (prefix string, useLivePrefix bool)
	Commands   []Command
}

type Register func() (*Runtime, error)

type Environment string

const (
	Docker Environment = "Docker"
)

var Environments = map[Environment]*Runtime{}

var Registers = map[Environment]Register{}

var ActiveEnvironment Environment

// Select one of installed environment
func SelectCurrentEnvironment() {
	if len(Environments) == 0 {
		fmt.Printf("There is no Docker envirnoment")
		os.Exit(0)
	} else if len(Environments) == 1 {
		for k, _ := range Environments {
			ActiveEnvironment = k
		}
	} else {
		// TODO 其他环境还未实现
	}
}

// 检测并初始化环境
func Initialize() {
	template.SetSurveyTemplate()

	extraCmds := DefaultCommands
	for k, v := range Registers {
		if env, err := TimeoutRegister(v, 5*time.Second); err == nil {
			Environments[k] = env
			extraCmds = append(extraCmds, env.Commands...)
		}
		time.Sleep(100 * time.Millisecond)
	}
	// clear entire line and ove cursor to beginning of the line 1 lines down.
	ansi.EraseInLine(2)
	ansi.CursorNextLine(1)

	for _, v := range extraCmds {
		var envs []Environment
		if len(v.Environments) == 0 {
			envs = []Environment{Docker}
		} else {
			envs = v.Environments
		}
		for _, e := range envs {
			if env, ok := Environments[e]; ok {
				c := env.Completer
				if c != nil {
					c.Arguments = append(c.Arguments, prompt.Suggest{Text: v.Text, Alias: v.Alias, Provider: v.Provider, Description: v.Description})
					c.Arguments = append(c.Arguments, v.ExtendArguments...)
					c.Options = append(c.Options, v.ExtendOptions...)
				}
			}
		}
	}

	for k := range Environments {
		c := Environments[k].Completer
		if c != nil {
			c.Setup(c)
		}
	}
}

func GetActive() *Runtime {
	return Environments[ActiveEnvironment]
}

func LivePrefix() (prefix string, useLivePrefix bool) {
	act := GetActive()
	if act.LivePrefix == nil {
		return defaultLivePrefix()
	} else {
		return act.LivePrefix()
	}
}

func defaultLivePrefix() (prefix string, useLivePrefix bool) {
	return GetActive().Prefix + " ", true
}

func GetActiveCompleter() prompt.Completer {
	return func(document prompt.Document) []prompt.Suggest {
		return GetActive().Completer.Complete(document)
	}
}

func GetActiveExecutor() prompt.Executor {
	return func(cmd string) {
		GetActive().Executor(cmd)
	}
}

func TimeoutRegister(fn Register, timeout time.Duration) (*Runtime, error) {
	ctx := context.Background()
	done := make(chan *Runtime, 1)
	err := make(chan error, 1)

	go func(ctx context.Context) {
		if r, e := fn(); e != nil {
			err <- e
		} else {
			done <- r
		}
	}(ctx)

	select {
	case r := <-done:
		return r, nil
	case e := <-err:
		return nil, e
	case <-time.After(timeout):
		return nil, errors.New("execution timeout")
	}
}
