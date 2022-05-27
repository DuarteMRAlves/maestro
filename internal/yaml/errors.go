package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"regexp"
	"sort"
	"strings"
)

type missingRequiredField struct {
	Field string
}

func (err *missingRequiredField) Error() string {
	return fmt.Sprintf("missing required field '%s'", err.Field)
}

type unknownFields struct {
	Fields []string
}

func (err *unknownFields) Error() string {
	return fmt.Sprintf("unknown fields '%s'", strings.Join(err.Fields, ","))
}

func typeErrorToError(typeErr *yaml.TypeError) error {
	var unkFields []string
	unkRegex := regexp.MustCompile(
		`line \d+: field (?P<field>[\w\W_]+) not found in type [\w\W_.]+`,
	)
	for _, errMsg := range typeErr.Errors {
		match := unkRegex.FindStringSubmatch(errMsg)
		if len(match) > 0 {
			unkFields = append(unkFields, match[unkRegex.SubexpIndex("field")])
		}
	}
	if len(unkFields) > 0 {
		sort.Strings(unkFields)
		return &unknownFields{Fields: unkFields}
	}
	return typeErr
}
