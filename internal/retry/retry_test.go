package retry

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestExponentialBackoff_Next(t *testing.T) {
	tests := map[string]struct {
		strat    *ExponentialBackoff
		expected []time.Duration
	}{
		"default params": {
			strat: &ExponentialBackoff{},
			expected: []time.Duration{
				defaultInitBackoff,
				defaultInitBackoff * defaultFact,
				defaultInitBackoff * defaultFact * defaultFact,
			},
		},
		"custom params": {
			strat:    NewExponentialBackoff(2, 3),
			expected: []time.Duration{2, 2 * 3, 2 * 3 * 3},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var actual []time.Duration
			for i := 0; i < len(tc.expected); i++ {
				actual = append(actual, tc.strat.Next())
			}
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Fatalf("backoff mismatch:\n%s", diff)
			}
		})
	}
}
