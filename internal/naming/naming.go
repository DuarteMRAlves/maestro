package naming

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

func IsValidAssetName(name apitypes.AssetName) bool {
	return IsValidName(string(name))
}

func IsValidStageName(name apitypes.StageName) bool {
	return IsValidName(string(name))
}

func IsValidOrchestrationName(name apitypes.OrchestrationName) bool {
	return IsValidName(string(name))
}

func IsValidName(name string) bool {
	return nameRegExp.MatchString(name)
}
