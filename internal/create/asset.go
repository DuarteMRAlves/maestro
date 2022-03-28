package create

import (
	"errors"
	"fmt"
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

type assetAlreadyExists struct{ name string }

func (err *assetAlreadyExists) Error() string {
	return fmt.Sprintf("asset '%s' already exists", err.name)
}

func Asset(storage AssetStorage) func(
	internal.AssetName,
	internal.Image,
) error {
	return func(name internal.AssetName, image internal.Image) error {
		if name.IsEmpty() {
			return EmptyAssetName
		}
		if image.IsEmpty() {
			return EmptyImageName
		}
		// Expect key not found
		_, err := storage.Load(name)
		if err == nil {
			return &assetAlreadyExists{name: name.Unwrap()}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}
		asset := internal.NewAsset(name, image)
		return storage.Save(asset)
	}
}
