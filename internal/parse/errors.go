package parse

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"regexp"
	"sort"
	"strings"
)

type MissingRequiredField struct {
	Field string
}

func (err *MissingRequiredField) Error() string {
	return fmt.Sprintf("missing required field '%s'", err.Field)
}

type UnknownFields struct {
	Fields []string
}

func (err *UnknownFields) Error() string {
	return fmt.Sprintf("unknown fields '%s'", strings.Join(err.Fields, ","))
}

func typeErrorToError(typeErr *yaml.TypeError) error {
	var unknownFields []string
	unknownRegex := regexp.MustCompile(
		`line \d+: field (?P<field>[\w\W_]+) not found in type [\w\W_.]+`,
	)
	for _, errMsg := range typeErr.Errors {
		unknownMatch := unknownRegex.FindStringSubmatch(errMsg)
		if len(unknownMatch) > 0 {
			unknownFields = append(
				unknownFields, unknownMatch[unknownRegex.SubexpIndex("field")],
			)
		}
	}
	if len(unknownFields) > 0 {
		sort.Strings(unknownFields)
		return &UnknownFields{Fields: unknownFields}
	}
	return typeErr
}
