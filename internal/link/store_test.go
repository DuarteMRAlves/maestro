package link

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/test"
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
				test.IsTrue(t, ok, "type assertion failed for store")

				err := st.Create(cfg)
				test.IsNil(t, err, "create error")
				test.DeepEqual(t, 1, lenLinks(st), "store size")
				stored, ok := st.links.Load(cfg.Name)
				test.IsTrue(t, ok, "link exists")
				s, ok := stored.(*Link)
				test.IsTrue(t, ok, "link type assertion failed")
				test.DeepEqual(t, cfg.Name, s.Name, "correct name")
				test.DeepEqual(
					t,
					cfg.SourceStage,
					s.SourceStage,
					"correct source stage")
				test.DeepEqual(
					t,
					cfg.SourceField,
					s.SourceField,
					"correct source field")
				test.DeepEqual(
					t,
					cfg.TargetStage,
					s.TargetStage,
					"correct target stage")
				test.DeepEqual(
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
