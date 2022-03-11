package domain

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewAssetName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"largeNameWithLetters", true},
		{"125largeName1565WithNumbers123", true},
		{"name-with-dashes-123-and-number-021", true},
		{"name/with-dash/and_underscore:tag-123", true},
		{"a", true},
		{"", true},
		{"name--double-dash", false},
		{"/back-slash-start", false},
		{"colon-end:", false},
		{"under-score-and-back-slash_/image", false},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				created, err := NewAssetName(test.name)
				if test.valid {
					assert.NilError(t, err, "create error")
					assert.Equal(t, test.name, created.Unwrap())
				} else {
					assert.Assert(t, created.IsEmpty())
					assert.Assert(t, errdefs.IsInvalidArgument(err))
					errMsg := fmt.Sprintf("invalid name '%v'", test.name)
					assert.Error(t, err, errMsg)
				}
			},
		)
	}
}
