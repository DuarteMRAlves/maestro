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
	emptyAssetName = errors.New("empty asset name")
	emptyImageName = errors.New("empty image name")
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
			return emptyAssetName
		}
		if image.IsEmpty() {
			return emptyImageName
		}
		// Expect key not found
		_, err := storage.Load(name)
		if err == nil {
			return &assetAlreadyExists{name: name.Unwrap()}
		}
		var nf interface{ NotFound() }
		if !errors.As(err, &nf) {
			return err
		}
		asset := internal.NewAsset(name, image)
		return storage.Save(asset)
	}
}
