package naming

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

func IsValidAssetName(name api.AssetName) bool {
	return IsValidName(string(name))
}

func IsValidStageName(name api.StageName) bool {
	return IsValidName(string(name))
}

func IsValidLinkName(name api.LinkName) bool {
	return IsValidName(string(name))
}

func IsValidOrchestrationName(name api.OrchestrationName) bool {
	return IsValidName(string(name))
}

func IsValidName(name string) bool {
	return nameRegExp.MatchString(name)
}
