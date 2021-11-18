package blueprint

import (
	"fmt"
)

type Link struct {
	SourceStage string
	SourceField string
	TargetStage string
	TargetField string
}

func (l *Link) Clone() *Link {
	return &Link{
		SourceStage: l.SourceStage,
		SourceField: l.SourceField,
		TargetStage: l.TargetStage,
		TargetField: l.TargetField,
	}
}

func (l *Link) String() string {
	return fmt.Sprintf(
		"Link{SourceStage:%v,SourceField:'%v',TargetStage:%v,TargetField:'%v'",
		l.SourceStage,
		l.SourceField,
		l.TargetStage,
		l.TargetField)
}
