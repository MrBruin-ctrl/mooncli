package suggestions

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

type GenericCompleter struct {
	Completer
	SuggestionProviders Repository
	Arguments           []prompt.Suggest
	Options             []prompt.Suggest
	Setup               func(*GenericCompleter)
}

func NewGenericCompleter(arguments []prompt.Suggest, options []prompt.Suggest, setup func(*GenericCompleter)) *GenericCompleter {
	return &GenericCompleter{
		SuggestionProviders: Repository{privateProviders: make(map[string]Provider)},
		Arguments:           arguments,
		Options:             options,
		Setup:               setup,
	}
}

func (g GenericCompleter) Complete(doc prompt.Document) []prompt.Suggest {
	args := strings.Split(doc.TextBeforeCursor(), " ")
	currentArg := doc.GetWordBeforeCursor()

	if doc.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	// 如果包含管道符，返回空
	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}
	// 如果光标前字符为'-'，则提供option补全
	if strings.HasPrefix(currentArg, "-") {
		return g.optionCompleter(args, currentArg)
	}

	commandArgs, isOptionValue, provider := g.excludeOptions(args)
	if isOptionValue {
		return g.optionValueCompleter(args, currentArg, provider)
	} else {
		return g.argumentCompleter(commandArgs, currentArg)
	}
}

func (g GenericCompleter) argumentCompleter(args []string, currentArg string) []prompt.Suggest {
	l := len(args)
	provider := Argument
	if l > 1 {
		set := g.SuggestionProviders.Provide(Argument, args[:l-1]...) // should return single record
		if len(set) == 0 {
			return []prompt.Suggest{}
		}
		provider = set[0].Provider
	}
	return prompt.FilterHasPrefix(g.SuggestionProviders.Provide(provider, args...), currentArg, true)
}

func (g GenericCompleter) optionCompleter(args []string, currentArg string) []prompt.Suggest {
	return prompt.FilterContains(g.SuggestionProviders.Provide(Option, args[:len(args)-1]...), currentArg, false)
}

func (g GenericCompleter) optionValueCompleter(args []string, currentArg string, provider string) []prompt.Suggest {
	return prompt.FilterHasPrefix(g.SuggestionProviders.Provide(provider, args...), currentArg, true)
}

func (g GenericCompleter) excludeOptions(args []string) ([]string, bool, string) {
	l := len(args)
	if l <= 1 {
		return args, false, ""
	}
	filtered := make([]string, 0, l)
	var (
		optionValueFlag     bool
		optionValueProvider string
	)
	optionSet := g.SuggestionProviders.Provide(Option, args[:l-2]...)
	for i := 0; i < len(args)-1; i++ {
		if optionValueFlag {
			optionValueFlag = false
			continue
		}
		if (i == 0 && args[0] == "") || (args[i] == "" && args[i-1] == "") {
			continue
		}

		if strings.HasPrefix(args[i], "-") {
			optionValueFlag = true
			if strings.Contains(args[i], "=") {
				optionValueFlag = false
			}
			for _, s := range optionSet {
				if strings.EqualFold(s.Text, args[i]) || strings.EqualFold(s.Alias, args[i]) {
					optionValueProvider = s.Provider
					optionValueFlag = s.Provider != "bool"
				}
			}
			continue
		}
		filtered = append(filtered, args[i])
	}
	filtered = append(filtered, args[l-1])
	return filtered, optionValueFlag, optionValueProvider
}
