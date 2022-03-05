package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
)

func RequestToResult(req domain.CreateAssetRequest) domain.AssetResult {
	name, err := domain.NewAssetName(req.Name)
	if err != nil {
		return domain.ErrAsset(err)
	}
	if !req.Image.Present() {
		return domain.SomeAsset(domain.NewAssetWithoutImage(name))
	}
	img, err := domain.NewImage(req.Image.Unwrap())
	if err != nil {
		return domain.ErrAsset(err)
	}
	return domain.SomeAsset(domain.NewAssetWithImage(name, img))
}

func ResultToResponse(res domain.AssetResult) domain.CreateAssetResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return domain.CreateAssetResponse{Err: errOpt}
}
