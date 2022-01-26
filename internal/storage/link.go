package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
)

type Link struct {
	name        api.LinkName
	sourceStage api.StageName
	sourceField string
	targetStage api.StageName
	targetField string
}

func NewLink(
	name api.LinkName,
	sourceStage api.StageName,
	sourceField string,
	targetStage api.StageName,
	targetField string,
) *Link {
	return &Link{
		name:        name,
		sourceStage: sourceStage,
		sourceField: sourceField,
		targetStage: targetStage,
		targetField: targetField,
	}
}

func (l *Link) Name() api.LinkName {
	return l.name
}

func (l *Link) SourceStage() api.StageName {
	return l.sourceStage
}

func (l *Link) SourceField() string {
	return l.sourceField
}

func (l *Link) TargetStage() api.StageName {
	return l.targetStage
}

func (l *Link) TargetField() string {
	return l.targetField
}

func (l *Link) Clone() *Link {
	return &Link{
		name:        l.name,
		sourceStage: l.sourceStage,
		sourceField: l.sourceField,
		targetStage: l.targetStage,
		targetField: l.targetField,
	}
}

func (l *Link) ToApi() *api.Link {
	return &api.Link{
		Name:        l.name,
		SourceStage: l.sourceStage,
		SourceField: l.sourceField,
		TargetStage: l.targetStage,
		TargetField: l.targetField,
	}
}

func (l *Link) String() string {
	return fmt.Sprintf(
		"Link{Name:%v,SourceStage:%v,SourceField:%v,TargetStage:%v,TargetField:%v",
		l.name,
		l.sourceStage,
		l.sourceField,
		l.targetStage,
		l.targetField,
	)
}
