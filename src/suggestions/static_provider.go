package suggestions

import (
	"github.com/Mrbruin-ctrl/moon-cli/src/internal/util"
	"github.com/c-bata/go-prompt"
	"strings"
)

func BuildFixedSelectionProvider(selections ...string) Provider {
	var ret []prompt.Suggest
	for _, v := range selections {
		ret = append(ret, prompt.Suggest{Text: v})
	}
	return func(args ...string) []prompt.Suggest {
		return ret
	}
}

func BuildStaticCompletionProvider(options []prompt.Suggest, fn SuggestionFilter) Provider {
	// setup index
	suggestions := make(map[string][]prompt.Suggest)
	for _, v := range options {
		if _, ok := suggestions[v.Argument]; !ok {
			suggestions[v.Argument] = []prompt.Suggest{}
		}
		suggestions[v.Argument] = append(suggestions[v.Argument], v)
	}
	// build provider
	return func(args ...string) []prompt.Suggest {
		return fn(args, suggestions)
	}
}

func BestEffortFilter(args []string, suggestions map[string][]prompt.Suggest) []prompt.Suggest {
	cmd := strings.ToLower(strings.Join(args, " "))
	var sug []prompt.Suggest
	var arg string
	for key := range suggestions {
		if strings.HasPrefix(cmd, key) {
			if len(key) >= len(arg) {
				arg = key
				sug = suggestions[key]
			}
		}
	}
	return sug
}

func LengthFilter(args []string, suggestions map[string][]prompt.Suggest) []prompt.Suggest {
	l := len(args)
	s := suggestions[strings.ToLower(strings.Join(args[:l-1], " "))]
	var ret []prompt.Suggest
	for _, v := range s {
		if util.HasPrefixIgnoreCase(v.Text, args[l-1]) || util.HasPrefixIgnoreCase(v.Alias, args[l-1]) {
			ret = append(ret, v)
		}
	}
	return ret
}
