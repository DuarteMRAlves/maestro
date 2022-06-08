package compiled

import "regexp"

var nameRegExp, _ = regexp.Compile(`^/?[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

func validateResourceName(name string) bool {
	return nameRegExp.MatchString(name)
}
