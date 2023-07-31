package environments

import "github.com/Mrbruin-ctrl/moon-cli/src/suggestions"

type RuntimeCompleter struct {
	*suggestions.GenericCompleter
	*Runtime
}
