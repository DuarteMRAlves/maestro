package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

func CreateAsset(storage AssetStorage) func(AssetRequest) AssetResponse {
	return func(req AssetRequest) AssetResponse {
		res := requestToAsset(req)
		res = domain.BindAsset(newVerifyDuplicateFn(storage))(res)
		res = domain.BindAsset(storage.Save)(res)
		return assetToResponse(res)
	}
}

func requestToAsset(req AssetRequest) domain.AssetResult {
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

func newVerifyDuplicateFn(verifier AssetExistsVerifier) func(domain.Asset) domain.AssetResult {
	return func(a domain.Asset) domain.AssetResult {
		if verifier.Verify(a.Name()) {
			err := errdefs.AlreadyExistsWithMsg(
				"asset '%v' already exists",
				a.Name().Unwrap(),
			)
			return domain.ErrAsset(err)
		}
		return domain.SomeAsset(a)
	}
}

func assetToResponse(res domain.AssetResult) AssetResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return AssetResponse{Err: errOpt}
}
