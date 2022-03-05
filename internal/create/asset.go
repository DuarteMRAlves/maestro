package create

import "github.com/DuarteMRAlves/maestro/internal/domain"

func RequestToResult(req AssetRequest) domain.AssetResult {
	name, err := domain.NewAssetName(req.Name)
	if err != nil {
		return domain.ErrAsset(err)
	}
	imgOpt := domain.NewEmptyImage()
	if req.Image.Present() {
		img, err := domain.NewImage(req.Image.Unwrap())
		if err != nil {
			return domain.ErrAsset(err)
		}
		imgOpt = domain.NewPresentImage(img)
	}
	return domain.SomeAsset(domain.NewAsset(name, imgOpt))
}

func ResultToResponse(res domain.AssetResult) AssetResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return AssetResponse{Err: errOpt}
}
