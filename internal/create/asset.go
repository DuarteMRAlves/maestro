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

var EmptyAssetName = fmt.Errorf("empty asset name")

func Asset(storage AssetStorage) func(AssetRequest) error {
	return func(req AssetRequest) error {
		name, err := internal.NewAssetName(req.Name)
		if err != nil {
			return err
		}
		if name.IsEmpty() {
			return EmptyAssetName
		}
		imgOpt := internal.NewEmptyImage()
		if req.Image.Present() {
			img := internal.NewImage(req.Image.Unwrap())
			imgOpt = internal.NewPresentImage(img)
		}
		// Expect key not found
		_, err = storage.Load(name)
		if err == nil {
			return &internal.AlreadyExists{Type: "asset", Ident: name.Unwrap()}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}
		asset := internal.NewAsset(name, imgOpt)
		return storage.Save(asset)
	}
}
