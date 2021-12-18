package stage

import (
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	stageName    = "stage-name"
	stageAsset   = "asset-name"
	stageService = "ServiceName"
	stageMethod  = "MethodName"
	stageAddress = "Address"
)

func TestStore_Create(t *testing.T) {
	tests := []struct {
		name   string
		config *Stage
	}{
		{
			name: "non default params",
			config: &Stage{
				Name:    stageName,
				Asset:   stageAsset,
				Service: stageService,
				Method:  stageMethod,
				Address: stageAddress,
			},
		},
		{
			name: "default params",
			config: &Stage{
				Name:    "",
				Asset:   "",
				Service: "",
				Method:  "",
				Address: "",
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
				assert.Equal(t, 1, lenStages(st), "store size")
				stored, ok := st.stages.Load(cfg.Name)
				assert.Assert(t, ok, "stage exists")
				s, ok := stored.(*Stage)
				assert.Assert(t, ok, "stage type assertion failed")
				assert.Equal(t, cfg.Name, s.Name, "correct name")
				assert.Equal(t, cfg.Asset, s.Asset, "correct asset")
				assert.Equal(t, cfg.Service, s.Service, "correct service")
				assert.Equal(t, cfg.Method, s.Method, "correct method")
			})
	}
}

func lenStages(st *store) int {
	count := 0
	st.stages.Range(
		func(key, value interface{}) bool {
			count += 1
			return true
		})
	return count
}

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name  string
		query *Stage
		// numbers to be stores
		stored []int
		// names of the expected stages
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
			query:    &Stage{Name: "some-name"},
			stored:   []int{},
			expected: []string{},
		},
		{
			name:     "one element stored, nil query",
			query:    nil,
			stored:   []int{0},
			expected: []string{testutil.StageNameForNum(0)},
		},
		{
			name:   "multiple elements stored, nil query",
			query:  nil,
			stored: []int{0, 1, 2},
			expected: []string{
				testutil.StageNameForNum(0),
				testutil.StageNameForNum(1),
				testutil.StageNameForNum(2),
			},
		},
		{
			name:     "multiple elements stored, matching name query",
			query:    &Stage{Name: testutil.StageNameForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.StageNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching name query",
			query:    &Stage{Name: "unknown-name"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching asset query",
			query:    &Stage{Asset: testutil.AssetNameForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.StageNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching asset query",
			query:    &Stage{Asset: "unknown-name"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching service query",
			query:    &Stage{Service: testutil.StageServiceForNum(1)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.StageNameForNum(1)},
		},
		{
			name:     "multiple elements stored, non-matching service query",
			query:    &Stage{Service: "unknown-service"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching method query",
			query:    &Stage{Method: testutil.StageMethodForNum(0)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.StageNameForNum(0)},
		},
		{
			name:     "multiple elements stored, non-matching method query",
			query:    &Stage{Method: "unknown-method"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name:     "multiple elements stored, matching address query",
			query:    &Stage{Address: testutil.StageAddressForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []string{testutil.StageNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching address query",
			query:    &Stage{Address: "unknown-address"},
			stored:   []int{0, 1, 2},
			expected: []string{},
		},
		{
			name: "multiple elements stored, exclusive query",
			query: &Stage{
				Asset:   testutil.AssetNameForNum(1),
				Address: testutil.StageAddressForNum(2),
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
					err := st.Create(stageForNum(n))
					assert.Assert(t, err, "create asset error")
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

func stageForNum(num int) *Stage {
	return &Stage{
		Name:    testutil.StageNameForNum(num),
		Asset:   testutil.AssetNameForNum(num),
		Service: testutil.StageServiceForNum(num),
		Method:  testutil.StageMethodForNum(num),
		Address: testutil.StageAddressForNum(num),
	}
}
