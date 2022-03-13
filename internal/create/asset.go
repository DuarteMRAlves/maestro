package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
)

type AssetSaver interface {
	Save(internal.Asset) error
}

type AssetLoader interface {
	Load(internal.AssetName) (internal.Asset, error)
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

var EmptyAssetName = fmt.Errorf("empty asset name")

func Asset(storage AssetStorage) func(AssetRequest) AssetResponse {
	return func(req AssetRequest) AssetResponse {
		name, err := internal.NewAssetName(req.Name)
		if err != nil {
			return AssetResponse{Err: domain.NewPresentError(err)}
		}
		if name.IsEmpty() {
			return AssetResponse{Err: domain.NewPresentError(EmptyAssetName)}
		}
		imgOpt := internal.NewEmptyImage()
		if req.Image.Present() {
			img := internal.NewImage(req.Image.Unwrap())
			imgOpt = internal.NewPresentImage(img)
		}
		// Expect key not found
		_, err = storage.Load(name)
		if err == nil {
			err := &internal.AlreadyExists{Type: "asset", Ident: name.Unwrap()}
			return AssetResponse{Err: domain.NewPresentError(err)}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return AssetResponse{Err: domain.NewPresentError(err)}
		}
		asset := internal.NewAsset(name, imgOpt)
		err = storage.Save(asset)
		errOpt := domain.NewEmptyError()
		if err != nil {
			errOpt = domain.NewPresentError(err)
		}
		return AssetResponse{Err: errOpt}
	}
}
