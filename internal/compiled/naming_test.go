package compiled

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidateResourceName(t *testing.T) {
	tests := map[string]struct {
		name  string
		valid bool
	}{
		"letters and numbers":     {name: "someName123", valid: true},
		"slashes":                 {name: "/org/repo/name", valid: true},
		"hifens":                  {name: "org-repo-name", valid: true},
		"colons":                  {name: "org:repo:name", valid: true},
		"mult special chars":      {name: "/org/repo/name-1:version", valid: true},
		"empty":                   {name: "", valid: true},
		"following special chars": {name: "some//name::1", valid: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := validateResourceName(tc.name)
			if diff := cmp.Diff(tc.valid, result); diff != "" {
				t.Fatalf("result mismatch:\n%s", diff)
			}
		})
	}
}
