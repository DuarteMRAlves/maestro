package orchestration

import (
	"fmt"
	"gotest.tools/v3/assert"
	"testing"
)

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"largeImageWithLetters", true},
		{"125largeImage1565WithNumbers123", true},
		{"image-with-dashes-123-and-number-021", true},
		{"image/with-dash/and_underscore:tag-123", true},
		{"a", true},
		{"", false},
		{"image--double-dash", false},
		{"/back-slash-start", false},
		{"colon-end:", false},
		{"under-score-and-back-slash_/image", false},
	}

	for _, inner := range tests {
		name, expected := inner.name, inner.expected
		testName := fmt.Sprintf("name='%v', expected=%v", name, expected)

		t.Run(
			testName,
			func(t *testing.T) {
				result := isValidName(name)
				assert.Equal(t, expected, result, "is valid name")
			},
		)
	}
}