package create

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

func ParseFiles(files []string) ([]*Resource, error) {
	resources := make([]*Resource, 0)

	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errdefs.InvalidArgumentWithError(err)
		}
		dataResources, err := ParseBytes(data)
		if err != nil {
			return nil, err
		}
		resources = append(resources, dataResources...)
	}

	return resources, nil
}

func ParseBytes(data []byte) ([]*Resource, error) {
	var err error

	reader := bytes.NewReader(data)

	dec := yaml.NewDecoder(reader)
	dec.SetStrict(true)

	resources := make([]*Resource, 0)

	for {
		r := &Resource{}
		if err = dec.Decode(&r); err != nil {
			break
		}
		resources = append(resources, r)
	}
	switch {
	case err == io.EOF:
		return resources, nil
	case errdefs.IsInvalidArgument(err):
		return nil, err
	default:
		typeErr, ok := err.(*yaml.TypeError)
		if ok {
			return nil, typeErrorToError(typeErr)
		}
		return nil, errdefs.InvalidArgumentWithMsg(
			"unmarshal resource: %v",
			err,
		)
	}
}

func typeErrorToError(typeErr *yaml.TypeError) error {
	unknownRegex := regexp.MustCompile(
		`line \d+: field (?P<field>[\w\W_]+) not found in type [\w\W_.]+`,
	)
	unknownFields := make([]string, 0)
	for _, errMsg := range typeErr.Errors {
		unknownMatch := unknownRegex.FindStringSubmatch(errMsg)
		if len(unknownMatch) > 0 {
			unknownFields = append(
				unknownFields,
				unknownMatch[unknownRegex.SubexpIndex("field")],
			)
		}
	}
	if len(unknownFields) > 0 {
		sort.Strings(unknownFields)
		return errdefs.InvalidArgumentWithMsg(
			"unknown spec fields: %v",
			strings.Join(unknownFields, ","),
		)
	}
	return errdefs.InvalidArgumentWithMsg(
		"unmarshal resource: %v",
		typeErr,
	)
}
