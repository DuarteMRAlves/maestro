package create

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type AssetSaver interface {
	Save(internal.Asset) internal.AssetResult
}

type AssetLoader interface {
	Load(internal.AssetName) internal.AssetResult
}

type AssetStorage interface {
	AssetSaver
	AssetLoader
}

type AssetRequest struct {
	Name  string
	Image domain.OptionalString
}

type AssetResponse struct {
	Err domain.OptionalError
}

func Asset(storage AssetStorage) func(AssetRequest) AssetResponse {
	return func(req AssetRequest) AssetResponse {
		res := requestToAsset(req)
		res = internal.BindAsset(newVerifyDuplicateFn(storage))(res)
		res = internal.BindAsset(storage.Save)(res)
		return assetToResponse(res)
	}
}

func requestToAsset(req AssetRequest) internal.AssetResult {
	name, err := internal.NewAssetName(req.Name)
	if err != nil {
		return internal.ErrAsset(err)
	}
	if name.IsEmpty() {
		err := errdefs.InvalidArgumentWithMsg("empty asset name")
		return internal.ErrAsset(err)
	}
	imgOpt := internal.NewEmptyImage()
	if req.Image.Present() {
		img := internal.NewImage(req.Image.Unwrap())
		imgOpt = internal.NewPresentImage(img)
	}
	return internal.SomeAsset(internal.NewAsset(name, imgOpt))
}

func newVerifyDuplicateFn(loader AssetLoader) func(internal.Asset) internal.AssetResult {
	return func(a internal.Asset) internal.AssetResult {
		res := loader.Load(a.Name())
		if res.IsError() {
			err := res.Error()
			if errdefs.IsNotFound(err) {
				return internal.SomeAsset(a)
			}
			return internal.ErrAsset(err)
		}
		err := errdefs.AlreadyExistsWithMsg(
			"asset '%v' already exists",
			a.Name().Unwrap(),
		)
		return internal.ErrAsset(err)
	}
}

func assetToResponse(res internal.AssetResult) AssetResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return AssetResponse{Err: errOpt}
}
