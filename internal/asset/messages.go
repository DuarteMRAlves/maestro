package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
)

func RequestToResult(req domain.CreateAssetRequest) domain.AssetResult {
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

func ResultToResponse(res domain.AssetResult) domain.CreateAssetResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return domain.CreateAssetResponse{Err: errOpt}
}
