package resources

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/test"
	"testing"
)

const (
	optional      = "optional"
	upperOptional = "upperOptional"
	required      = "required"
	upperRequired = "upperRequired"
)

type testResource struct {
	Optional      string
	UpperOptional string `yaml:"upper_optional"`
	Required      string `yaml:",required"`
	UpperRequired string `yaml:"upper_required,required"`
	unexported    string
}

var (
	inputAllFields = &Resource{
		Kind: "test",
		Spec: map[string]string{
			"Optional":       optional,
			"upper_optional": upperOptional,
			"Required":       required,
			"upper_required": upperRequired,
		},
	}
	expectedAllFields = &testResource{
		Optional:      optional,
		UpperOptional: upperOptional,
		Required:      required,
		UpperRequired: upperRequired,
	}
	inputRequiredFields = &Resource{
		Kind: "test",
		Spec: map[string]string{
			"Required":       required,
			"upper_required": upperRequired,
		},
	}
	expectedRequiredFields = &testResource{
		Required:      required,
		UpperRequired: upperRequired,
	}
	// Inputs with errors
	missingRequiredField = &Resource{
		Kind: "test",
		Spec: map[string]string{
			"Optional":       optional,
			"upper_optional": upperOptional,
			"Required":       required,
		},
	}
	unknownField = &Resource{
		Kind: "test",
		Spec: map[string]string{
			"Required":        required,
			"upper_required":  upperRequired,
			"unknown_field_1": "unknown1",
			"unknown_field_2": "unknown2",
		},
	}
)

func TestMarshalLinkResource_Correct(t *testing.T) {
	tests := []struct {
		input    *Resource
		expected *testResource
	}{
		{inputAllFields, expectedAllFields},
		{inputRequiredFields, expectedRequiredFields},
	}
	for _, inner := range tests {
		input, expected := inner.input, inner.expected
		testName := fmt.Sprintf("src=%v,expected=%v", input, expected)

		t.Run(
			testName,
			func(t *testing.T) {
				var resource testResource
				err := MarshalResource(&resource, input)
				test.IsNil(t, err, "error not nil")
				test.DeepEqual(t, expected, &resource, "resource differs")
			})
	}
}

func TestMarshalResource_Incorrect(t *testing.T) {
	intVar := 1
	tests := []struct {
		name     string
		dst      interface{}
		src      *Resource
		expected error
	}{
		{
			"dst is nil",
			nil,
			&Resource{},
			errdefs.InvalidArgumentWithMsg("'dst' is nil"),
		},
		{
			"src is nil",
			&testResource{},
			nil,
			errdefs.InvalidArgumentWithMsg("'src' is nil"),
		},
		{
			"dst not a pointer",
			testResource{},
			&Resource{},
			errdefs.InvalidArgumentWithMsg("dst must be a pointer"),
		},
		{
			"dst not a pointer to struct",
			&intVar,
			&Resource{},
			errdefs.InvalidArgumentWithMsg(
				"underlying dst object must be a struct"),
		},
		{
			"missing required field",
			&testResource{},
			missingRequiredField,
			errdefs.InvalidArgumentWithMsg("missing spec field upper_required"),
		},
		{
			"unknown field",
			&testResource{},
			unknownField,
			errdefs.InvalidArgumentWithMsg(
				"unknown spec fields: unknown_field_1,unknown_field_2"),
		},
	}
	for _, inner := range tests {
		testName := inner.name
		dst, src, expected := inner.dst, inner.src, inner.expected

		t.Run(
			testName,
			func(t *testing.T) {
				err := MarshalResource(dst, src)
				test.DeepEqual(t, expected, err, "error differs")
			})
	}
}
