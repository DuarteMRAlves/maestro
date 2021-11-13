package blueprint

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

type Link struct {
	SourceId    identifier.Id
	SourceField string
	TargetId    identifier.Id
	TargetField string
}

func (l *Link) Clone() *Link {
	return &Link{
		SourceId:    l.SourceId,
		SourceField: l.SourceField,
		TargetId:    l.TargetId,
		TargetField: l.TargetField,
	}
}

func (l *Link) String() string {
	return fmt.Sprintf(
		"Link{SourceId:%v,SourceField:'%v',TargetId:%v,TargetField:'%v'",
		l.SourceId,
		l.SourceField,
		l.TargetId,
		l.TargetField)
}
