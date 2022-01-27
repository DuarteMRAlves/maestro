package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
)

func UnmarshalCreateOrchestrationRequest(
	req *api.CreateOrchestrationRequest,
	pbReq *pb.CreateOrchestrationRequest,
) {
	req.Name = api.OrchestrationName(pbReq.Name)
}

func UnmarshalGetOrchestrationRequest(
	req *api.GetOrchestrationRequest,
	pbReq *pb.GetOrchestrationRequest,
) {
	req.Name = api.OrchestrationName(pbReq.Name)
	req.Phase = api.OrchestrationPhase(pbReq.Phase)
}

func UnmarshalCreateStageRequest(
	req *api.CreateStageRequest,
	pbReq *pb.CreateStageRequest,
) {
	req.Name = api.StageName(pbReq.Name)
	req.Service = pbReq.Service
	req.Rpc = pbReq.Rpc
	req.Address = pbReq.Address
	req.Host = pbReq.Host
	req.Port = pbReq.Port
	req.Asset = api.AssetName(pbReq.Asset)
}

func UnmarshalGetStageRequest(
	req *api.GetStageRequest,
	pbReq *pb.GetStageRequest,
) {
	req.Name = api.StageName(pbReq.Name)
	req.Phase = api.StagePhase(pbReq.Phase)
	req.Service = pbReq.Service
	req.Rpc = pbReq.Rpc
	req.Address = pbReq.Address
	req.Asset = api.AssetName(pbReq.Asset)
}

func UnmarshalCreateLinkRequest(
	req *api.CreateLinkRequest,
	pbReq *pb.CreateLinkRequest,
) {
	req.Name = api.LinkName(pbReq.Name)
	req.SourceStage = api.StageName(pbReq.SourceStage)
	req.SourceField = pbReq.SourceField
	req.TargetStage = api.StageName(pbReq.TargetStage)
	req.TargetField = pbReq.TargetField
}

func UnmarshalGetLinkRequest(
	req *api.GetLinkRequest,
	pbReq *pb.GetLinkRequest,
) {
	req.Name = api.LinkName(pbReq.Name)
	req.SourceStage = api.StageName(pbReq.SourceStage)
	req.SourceField = pbReq.SourceField
	req.TargetStage = api.StageName(pbReq.TargetStage)
	req.TargetField = pbReq.TargetField
}

func UnmarshalCreateAssetRequest(
	req *api.CreateAssetRequest,
	pbReq *pb.CreateAssetRequest,
) {
	req.Name = api.AssetName(pbReq.Name)
	req.Image = pbReq.Image
}

func UnmarshalGetAssetRequest(
	req *api.GetAssetRequest,
	pbReq *pb.GetAssetRequest,
) {
	req.Name = api.AssetName(pbReq.Name)
	req.Image = pbReq.Image
}
