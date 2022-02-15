package resources

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"reflect"
	"strings"
)

// validateInfo verifies if all the restrictions specified by the info tags are
// complied
func validateInfo(v interface{}) error {
	var (
		ok  bool
		err error
	)

	if ok, err = util.ArgNotNil(v, "dst"); !ok {
		return err
	}
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem()
	case reflect.Struct:
		// Do nothing, we keep the same value to later analyze as we already
		// have the struct.
	default:
		return errdefs.InvalidArgumentWithMsg(
			"invalid type: expected Ptr or Struct but got %v",
			value.Kind(),
		)
	}

	objType := value.Type()
	for i := 0; i < objType.NumField(); i++ {
		objTypeField := objType.Field(i)
		// Ignore unexported fields
		if !objTypeField.IsExported() {
			continue
		}
		fieldValue := value.Field(i)
		if err := validateField(objTypeField, fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func validateField(
	objTypeField reflect.StructField,
	fieldValue reflect.Value,
) error {
	tag, hasTag := objTypeField.Tag.Lookup("info")
	if hasTag {

		tagOpts := strings.Split(tag, ",")
		for _, opt := range tagOpts {
			switch opt {
			case "required":
				if fieldValue.IsZero() {
					return errdefs.InvalidArgumentWithMsg(
						"missing required field: '%v'",
						yamlName(objTypeField),
					)
				}
			}
		}
	}
	return nil
}

func yamlName(f reflect.StructField) string {
	tag, hasTag := f.Tag.Lookup("yaml")
	if hasTag {
		tagOpts := strings.Split(tag, ",")
		if tagOpts[0] != "" {
			return tagOpts[0]
		}
	}
	return f.Name
}
