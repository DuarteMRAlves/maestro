package pb

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
)

// MarshalAsset returns a new protobuf representations of the received asset.
func MarshalAsset(a *api.Asset) (*pb.Asset, error) {
	if ok, err := util.ArgNotNil(a, "a"); !ok {
		return nil, err
	}
	return &pb.Asset{Name: string(a.Name), Image: a.Image}, nil
}

// UnmarshalAsset returns a new asset from its protobuf representation.
func UnmarshalAsset(p *pb.Asset) (*api.Asset, error) {
	if ok, err := util.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &api.Asset{
		Name:  api.AssetName(p.Name),
		Image: p.Image,
	}, nil
}

// MarshalOrchestration returns a protobuf encoding of the given orchestration.
func MarshalOrchestration(o *api.Orchestration) (
	*pb.Orchestration,
	error,
) {
	if ok, err := util.ArgNotNil(o, "o"); !ok {
		return nil, err
	}
	links := make([]string, 0, len(o.Links))
	for _, l := range o.Links {
		links = append(links, string(l))
	}
	protoBp := &pb.Orchestration{
		Name:  string(o.Name),
		Phase: string(o.Phase),
		Links: links,
	}
	return protoBp, nil
}

// UnmarshalOrchestration returns an orchestration from the orchestration
// protobuf encoding.
func UnmarshalOrchestration(p *pb.Orchestration) (
	*api.Orchestration,
	error,
) {
	if ok, err := util.ArgNotNil(p, "p"); !ok {
		return nil, err
	}

	links := make([]api.LinkName, 0, len(p.Links))
	for _, l := range p.Links {
		links = append(links, api.LinkName(l))
	}

	return &api.Orchestration{
		Name:  api.OrchestrationName(p.Name),
		Phase: api.OrchestrationPhase(p.Phase),
		Links: links,
	}, nil
}

// MarshalStage creates a protobuf message representing the Stage from the Stage
// structure.
func MarshalStage(s *api.Stage) (*pb.Stage, error) {
	if ok, err := util.ArgNotNil(s, "s"); !ok {
		return nil, err
	}
	pbStage := &pb.Stage{
		Name:    string(s.Name),
		Phase:   string(s.Phase),
		Asset:   string(s.Asset),
		Service: s.Service,
		Rpc:     s.Rpc,
		Address: s.Address,
	}
	return pbStage, nil
}

// UnmarshalStage creates a Stage struct from a protobuf message representing
// the stage.
func UnmarshalStage(p *pb.Stage) (*api.Stage, error) {
	if ok, err := util.ArgNotNil(p, "p"); !ok {
		return nil, errdefs.InvalidArgumentWithError(err)
	}
	return &api.Stage{
		Name:    api.StageName(p.Name),
		Phase:   api.StagePhase(p.Phase),
		Asset:   api.AssetName(p.Asset),
		Service: p.Service,
		Rpc:     p.Rpc,
		Address: p.Address,
	}, nil
}

// MarshalLink returns a protobuf message for the given link.
func MarshalLink(l *api.Link) (*pb.Link, error) {
	if ok, err := util.ArgNotNil(l, "l"); !ok {
		return nil, err
	}
	pbLink := &pb.Link{
		Name:        string(l.Name),
		SourceStage: string(l.SourceStage),
		SourceField: l.SourceField,
		TargetStage: string(l.TargetStage),
		TargetField: l.TargetField,
	}
	return pbLink, nil
}

// UnmarshalLink returns the link represented by the given protobuf message.
func UnmarshalLink(p *pb.Link) (*api.Link, error) {
	if ok, err := util.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &api.Link{
		Name:        api.LinkName(p.Name),
		SourceStage: api.StageName(p.SourceStage),
		SourceField: p.SourceField,
		TargetStage: api.StageName(p.TargetStage),
		TargetField: p.TargetField,
	}, nil
}
