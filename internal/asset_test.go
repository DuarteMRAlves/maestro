package internal

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"reflect"
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
					if err != nil {
						t.Fatalf("error not nil in valid test: %s", err)
					}
					if diff := cmp.Diff(tc.name, created.Unwrap()); diff != "" {
						t.Fatalf("names mismatch:\n%s", diff)
					}
				} else {
					if !created.IsEmpty() {
						t.Fatalf("created is not empty: %s", created.Unwrap())
					}
					var invalidIdentifier *InvalidIdentifier
					if !errors.As(err, &invalidIdentifier) {
						errTyp := reflect.TypeOf(err)
						t.Fatalf(
							"error type mismatch: expected *InvalidIdentifier, got %s",
							errTyp,
						)
					}
					if diff := cmp.Diff(
						"asset",
						invalidIdentifier.Type,
					); diff != "" {
						t.Fatalf("error type mismatch:\n%s", diff)
					}
					if diff := cmp.Diff(
						tc.name,
						invalidIdentifier.Ident,
					); diff != "" {
						t.Fatalf("error identifier mismatch:\n%s", diff)
					}
				}
			},
		)
	}
}
