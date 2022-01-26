package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

func UnmarshalCreateOrchestrationRequest(
	req *api.CreateOrchestrationRequest,
	pbReq *pb.CreateOrchestrationRequest,
) {
	req.Name = apitypes.OrchestrationName(pbReq.Name)
	req.Links = make([]apitypes.LinkName, 0, len(pbReq.Links))
	for _, l := range pbReq.Links {
		req.Links = append(req.Links, apitypes.LinkName(l))
	}
}

func UnmarshalGetOrchestrationRequest(
	req *api.GetOrchestrationRequest,
	pbReq *pb.GetOrchestrationRequest,
) {
	req.Name = apitypes.OrchestrationName(pbReq.Name)
	req.Phase = apitypes.OrchestrationPhase(pbReq.Phase)
}

func UnmarshalCreateAssetRequest(
	req *api.CreateAssetRequest,
	pbReq *pb.CreateAssetRequest,
) {
	req.Name = apitypes.AssetName(pbReq.Name)
	req.Image = pbReq.Image
}

func UnmarshalGetAssetRequest(
	req *api.GetAssetRequest,
	pbReq *pb.GetAssetRequest,
) {
	req.Name = apitypes.AssetName(pbReq.Name)
	req.Image = pbReq.Image
}
