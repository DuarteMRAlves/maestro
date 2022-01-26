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
	req.Links = make([]api.LinkName, 0, len(pbReq.Links))
	for _, l := range pbReq.Links {
		req.Links = append(req.Links, api.LinkName(l))
	}
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
