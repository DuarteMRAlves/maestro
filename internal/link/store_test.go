package link

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	linkName        = "link-name"
	linkSourceStage = "linkSourceStage"
	linkSourceField = "linkSourceField"
	linkTargetStage = "linkTargetStage"
	linkTargetField = "linkTargetField"
)

func TestStore_Create(t *testing.T) {
	tests := []struct {
		name   string
		config *Link
	}{
		{
			name: "non default parameters",
			config: &Link{
				Name:        linkName,
				SourceStage: linkSourceStage,
				SourceField: linkSourceField,
				TargetStage: linkTargetStage,
				TargetField: linkTargetField,
			},
		},
		{
			name: "default parameters",
			config: &Link{
				Name:        "",
				SourceStage: "",
				SourceField: "",
				TargetStage: "",
				TargetField: "",
			},
		},
	}

	for _, test := range tests {

		t.Run(
			test.name,
			func(t *testing.T) {
				cfg := test.config

				st, ok := NewStore().(*store)
				assert.Assert(t, ok, "type assertion failed for store")

				err := st.Create(cfg)
				assert.NilError(t, err, "create error")
				assert.Equal(t, 1, lenLinks(st), "store size")
				stored, ok := st.links.Load(cfg.Name)
				assert.Assert(t, ok, "link exists")
				s, ok := stored.(*Link)
				assert.Assert(t, ok, "link type assertion failed")
				assert.Equal(t, cfg.Name, s.Name, "correct name")
				assert.Equal(
					t,
					cfg.SourceStage,
					s.SourceStage,
					"correct source stage")
				assert.Equal(
					t,
					cfg.SourceField,
					s.SourceField,
					"correct source field")
				assert.Equal(
					t,
					cfg.TargetStage,
					s.TargetStage,
					"correct target stage")
				assert.Equal(
					t,
					cfg.TargetField,
					s.TargetField,
					"correct target field")
			})
	}
}

func lenLinks(st *store) int {
	count := 0
	st.links.Range(
		func(key, value interface{}) bool {
			count += 1
			return true
		})
	return count
}

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name  string
		query *apitypes.Link
		// numbers to be stored
		stored []int
		// names of the expected Links
		expected []string
	}{
		{
			name:     "zero elements stored, nil query",
			query:    nil,
			stored:   []int{},
			expected: []string{},
		},
		{
			name:     "zero elements stored, some query",
			query:    &apitypes.Link{Name: "some-name"},
			stored:   []int{},
			expected: []string{},
		},
		{
			name:     "one element stored, nil query",
			query:    nil,
			stored:   []int{0},
			expected: []string{testutil.LinkNameForNum(0)},
		},
		{
			name:   "multiple elements stored, nil query",
			query:  nil,
			stored: []int{0, 1, 2},
			expected: []string{
				testutil.LinkNameForNum(0),
				testutil.LinkNameForNum(1),
				testutil.LinkNameForNum(2),
			},
		},
		{
			name:     "multiple elements stored, matching name query",
			query:    &apitypes.Link{Name: testutil.LinkNameForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.LinkNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching name query",
			query:    &apitypes.Link{Name: "unknown-name"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching source stage query",
			query:    &apitypes.Link{SourceStage: testutil.LinkSourceStageForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.LinkNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching source stage query",
			query:    &apitypes.Link{SourceStage: "unknown-stage"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching source field query",
			query:    &apitypes.Link{SourceField: testutil.LinkSourceFieldForNum(1)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.LinkNameForNum(1)},
		},
		{
			name:     "multiple elements stored, non-matching source field query",
			query:    &apitypes.Link{SourceField: "UnknownField"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching target stage query",
			query:    &apitypes.Link{TargetStage: testutil.LinkTargetStageForNum(0)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.LinkNameForNum(0)},
		},
		{
			name:     "multiple elements stored, non-matching target stage query",
			query:    &apitypes.Link{TargetStage: "unknown-stage"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching target field query",
			query:    &apitypes.Link{TargetField: testutil.LinkTargetFieldForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.LinkNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching target field query",
			query:    &apitypes.Link{TargetField: "UnknownField"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name: "multiple elements stored, exclusive query",
			query: &apitypes.Link{
				SourceStage: testutil.LinkSourceStageForNum(1),
				TargetStage: testutil.LinkTargetStageForNum(2),
			},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				st := NewStore()

				for _, n := range test.stored {
					err := st.Create(linkForNum(n))
					assert.NilError(t, err, "create Link error")
				}

				received := st.Get(test.query)
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[string]bool, 0)
				for _, e := range test.expected {
					seen[e] = false
				}

				for _, r := range received {
					alreadySeen, exists := seen[r.Name]
					assert.Assert(t, exists, "element should be expected")
					// Elements can't be seen twice
					assert.Assert(t, !alreadySeen, "element already seen")
					seen[r.Name] = true
				}

				for _, e := range test.expected {
					// All elements should be seen
					assert.Assert(t, seen[e], "element not seen")
				}
			})
	}
}

func linkForNum(num int) *Link {
	return &Link{
		Name:        testutil.LinkNameForNum(num),
		SourceStage: testutil.LinkSourceStageForNum(num),
		SourceField: testutil.LinkSourceFieldForNum(num),
		TargetStage: testutil.LinkTargetStageForNum(num),
		TargetField: testutil.LinkTargetFieldForNum(num),
	}
}
