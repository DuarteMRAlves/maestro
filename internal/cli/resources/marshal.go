package resources

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"reflect"
	"sort"
	"strings"
)

func MarshalResource(dst interface{}, src *Resource) error {
	if err := validateMarshalResourceArgs(dst, src); err != nil {
		return err
	}

	ptrValue := reflect.ValueOf(dst)
	objValue := ptrValue.Elem()
	objType := objValue.Type()

	// Copy spec to later destroy
	spec := make(map[string]string, len(src.Spec))
	for k, v := range src.Spec {
		spec[k] = v
	}
	for i := 0; i < objType.NumField(); i++ {
		var (
			key           string
			fieldRequired bool
		)
		objTypeField := objType.Field(i)
		// Ignore unexported fields
		if !objTypeField.IsExported() {
			continue
		}

		objValueField := objValue.Field(i)

		err := parseTag(&key, &fieldRequired, objTypeField)
		if err != nil {
			return err
		}

		value, exists := spec[key]
		if exists {
			objValueField.SetString(value)
			delete(spec, key)
		} else if fieldRequired {
			return errdefs.InvalidArgumentWithMsg("missing spec field %v", key)
		}
	}
	// Keys that did not match to any optional
	// Raise error for unknown keys
	if len(spec) > 0 {
		keys := make([]string, 0, len(spec))
		for k := range spec {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		keysDisplay := strings.Join(keys, ",")
		return errdefs.InvalidArgumentWithMsg(
			"unknown spec fields: %v",
			keysDisplay)
	}
	return nil
}

func validateMarshalResourceArgs(dst interface{}, src *Resource) error {
	var (
		ok  bool
		err error
	)
	if ok, err = validate.ArgNotNil(dst, "dst"); !ok {
		return err
	}
	if ok, err = validate.ArgNotNil(src, "src"); !ok {
		return err
	}
	ptrValue := reflect.ValueOf(dst)

	ok, err = validate.ArgStatus(
		ptrValue.Kind() == reflect.Ptr,
		"dst must be a pointer")
	if !ok {
		return err
	}
	objValue := ptrValue.Elem()
	ok, err = validate.ArgStatus(
		objValue.Kind() == reflect.Struct && !ptrValue.IsNil(),
		"underlying dst object must be a struct")
	if !ok {
		return err
	}
	return nil
}

func parseTag(name *string, required *bool, f reflect.StructField) error {
	// Fill with default values
	*name = f.Name
	*required = false

	// Process tags
	tag, hasTag := f.Tag.Lookup("yaml")
	if hasTag {
		tagOpts := strings.Split(tag, ",")
		if tagOpts[0] != "" {
			*name = tagOpts[0]
		}
		if len(tagOpts) > 1 {
			for _, opt := range tagOpts[1:] {
				switch opt {
				case "required":
					*required = true
					break
				default:
					return errdefs.InternalWithMsg("unknown tag %v", opt)
				}
			}
		}
	}
	return nil
}
