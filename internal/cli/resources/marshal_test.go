package resources

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
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
		testName := fmt.Sprintf("src=%v,errMsg=%v", input, expected)

		t.Run(
			testName,
			func(t *testing.T) {
				var resource testResource
				err := MarshalResource(&resource, input)
				assert.NilError(t, err, "error not nil")
				assert.Equal(t, expected.Optional, resource.Optional)
				assert.Equal(t, expected.UpperOptional, resource.UpperOptional)
				assert.Equal(t, expected.Required, resource.Required)
				assert.Equal(t, expected.UpperRequired, resource.UpperRequired)
				assert.Equal(t, expected.unexported, resource.unexported)
			})
	}
}

func TestMarshalResource_Incorrect(t *testing.T) {
	intVar := 1
	tests := []struct {
		name   string
		dst    interface{}
		src    *Resource
		errMsg string
	}{
		{"dst is nil", nil, &Resource{}, "'dst' is nil"},
		{"src is nil", &testResource{}, nil, "'src' is nil"},
		{
			"dst not a pointer",
			testResource{},
			&Resource{},
			"dst must be a pointer",
		},
		{
			"dst not a pointer to struct",
			&intVar,
			&Resource{},
			"underlying dst object must be a struct",
		},
		{
			"missing required field",
			&testResource{},
			missingRequiredField,
			"missing spec field upper_required",
		},
		{
			"unknown field",
			&testResource{},
			unknownField,
			"unknown spec fields: unknown_field_1,unknown_field_2",
		},
	}
	for _, inner := range tests {
		testName := inner.name
		dst, src, errMsg := inner.dst, inner.src, inner.errMsg

		t.Run(
			testName,
			func(t *testing.T) {
				err := MarshalResource(dst, src)
				assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
				assert.ErrorContains(t, err, errMsg)
			})
	}
}
