package suggestions

import (
	"github.com/c-bata/go-prompt"
)

type Completer interface {
	Complete(doc prompt.Document) []prompt.Suggest
}

type Provider func(...string) []prompt.Suggest

type SuggestionFilter func([]string, map[string][]prompt.Suggest) []prompt.Suggest

type Repository struct {
	privateProviders map[string]Provider
}

func (r *Repository) Add(id string, provider Provider) {
	r.privateProviders[id] = provider
}

func (r *Repository) Provide(id string, args ...string) []prompt.Suggest {
	if v, ok := r.privateProviders[id]; ok {
		return v(args...)
	} else if v, ok := sharedProviders[id]; ok {
		return v(args...)
	} else {
		return []prompt.Suggest{}
	}
}

//共享provider
var sharedProviders map[string]Provider

func init() {
	sharedProviders = make(map[string]Provider)
	RegisterSharedProvider(Path, providePathSuggestion)
	RegisterSharedProvider(Output, BuildFixedSelectionProvider("yaml", "json", "table", "short"))
	RegisterSharedProvider(Loglevel, BuildFixedSelectionProvider("trace", "debug", "info", "warn", "error", "fatal"))
}

func RegisterSharedProvider(id string, provider Provider) {
	sharedProviders[id] = provider
}

const (
	Option   = "option"
	Argument = "argument"
)
