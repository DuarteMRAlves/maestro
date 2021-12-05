package link

import (
	"fmt"
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
	tests := []*Link{
		{
			Name:        linkName,
			SourceStage: linkSourceStage,
			SourceField: linkSourceField,
			TargetStage: linkTargetStage,
			TargetField: linkTargetField,
		},
		{
			Name:        "",
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	}

	for _, cfg := range tests {
		testName := fmt.Sprintf("cfg=%v", cfg)

		t.Run(
			testName,
			func(t *testing.T) {
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
