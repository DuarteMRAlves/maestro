package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
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

var (
	EmptyAssetName = errors.New("empty asset name")
	EmptyImageName = errors.New("empty image name")
)

func Asset(storage AssetStorage) func(
	internal.AssetName,
	internal.OptionalImage,
) error {
	return func(name internal.AssetName, image internal.OptionalImage) error {
		if name.IsEmpty() {
			return EmptyAssetName
		}
		if image.Present() && image.Unwrap().IsEmpty() {
			return EmptyImageName
		}
		// Expect key not found
		_, err := storage.Load(name)
		if err == nil {
			return &internal.AlreadyExists{Type: "asset", Ident: name.Unwrap()}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}
		asset := internal.NewAsset(name, image)
		return storage.Save(asset)
	}
}
