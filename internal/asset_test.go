package internal

import (
	"errors"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewAssetName(t *testing.T) {
	tests := map[string]struct {
		name  string
		valid bool
	}{
		"letters":        {name: "simpleName", valid: true},
		"numbers":        {name: "125name1565num123", valid: true},
		"eq punct":       {name: "dashes-123-and-021", valid: true},
		"diff punct":     {name: "name-dash/and_under:tag-123", valid: true},
		"empty":          {name: "", valid: true},
		"start punct":    {name: "/back-slash-start", valid: false},
		"end punct":      {name: "colon-end:", valid: false},
		"two eq punct":   {name: "name--double", valid: false},
		"two diff punct": {name: "name_/image", valid: false},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				created, err := NewAssetName(tc.name)
				if tc.valid {
					assert.NilError(t, err, "create error")
					assert.Equal(t, tc.name, created.Unwrap())
				} else {
					assert.Assert(t, created.IsEmpty())
					var invalidIdentifier *InvalidIdentifier
					assert.Assert(t, errors.As(err, &invalidIdentifier))
					assert.Equal(t, "asset", invalidIdentifier.Type)
					assert.Equal(t, tc.name, invalidIdentifier.Ident)
				}
			},
		)
	}
}
