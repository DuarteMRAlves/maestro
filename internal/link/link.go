package link

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

type Link struct {
	Name        string
	SourceStage string
	SourceField string
	TargetStage string
	TargetField string
}

func (l *Link) Clone() *Link {
	return &Link{
		Name:        l.Name,
		SourceStage: l.SourceStage,
		SourceField: l.SourceField,
		TargetStage: l.TargetStage,
		TargetField: l.TargetField,
	}
}

func (l *Link) ToApi() *apitypes.Link {
	return &apitypes.Link{
		Name:        l.Name,
		SourceStage: l.SourceStage,
		SourceField: l.SourceField,
		TargetStage: l.TargetStage,
		TargetField: l.TargetField,
	}
}

func (l *Link) String() string {
	return fmt.Sprintf(
		"Link{Name:%v,SourceStage:%v,SourceField:%v,TargetStage:%v,TargetField:%v",
		l.Name,
		l.SourceStage,
		l.SourceField,
		l.TargetStage,
		l.TargetField)
}
