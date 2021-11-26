package link

import (
	"fmt"
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
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
				testing2.IsTrue(t, ok, "type assertion failed for store")

				err := st.Create(cfg)
				testing2.IsNil(t, err, "create error")
				testing2.DeepEqual(t, 1, lenLinks(st), "store size")
				stored, ok := st.links.Load(cfg.Name)
				testing2.IsTrue(t, ok, "link exists")
				s, ok := stored.(*Link)
				testing2.IsTrue(t, ok, "link type assertion failed")
				testing2.DeepEqual(t, cfg.Name, s.Name, "correct name")
				testing2.DeepEqual(
					t,
					cfg.SourceStage,
					s.SourceStage,
					"correct source stage")
				testing2.DeepEqual(
					t,
					cfg.SourceField,
					s.SourceField,
					"correct source field")
				testing2.DeepEqual(
					t,
					cfg.TargetStage,
					s.TargetStage,
					"correct target stage")
				testing2.DeepEqual(
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
