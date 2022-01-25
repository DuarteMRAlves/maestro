package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

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
