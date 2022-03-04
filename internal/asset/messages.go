package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/types"
)

func RequestToResult(req types.CreateAssetRequest) types.AssetResult {
	name, err := NewAssetName(req.Name)
	if err != nil {
		return NewErrResult(err)
	}
	if !req.Image.Present() {
		return NewResult(NewAssetWithoutImage(name))
	}
	img, err := NewImage(req.Image.Unwrap())
	if err != nil {
		return NewErrResult(err)
	}
	return NewResult(NewAssetWithImage(name, img))
}

func ResultToResponse(res types.AssetResult) types.CreateAssetResponse {
	var errOpt types.OptionalError
	if res.IsError() {
		errOpt = types.NewPresentError(res.Error())
	} else {
		errOpt = types.NewEmptyError()
	}
	return types.CreateAssetResponse{Err: errOpt}
}
