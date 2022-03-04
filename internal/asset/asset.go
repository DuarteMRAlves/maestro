package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/types"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type assetName string

func NewAssetName(name string) (types.AssetName, error) {
	if isValidResourceName(name) {
		return assetName(name), nil
	}
	return nil, errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
}

func isValidResourceName(name string) bool {
	return nameRegExp.MatchString(name)
}

func (a assetName) Unwrap() string {
	return string(a)
}

type image string

func NewImage(img string) (types.Image, error) {
	if len(img) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty image")
	}
	return image(img), nil
}

func (i image) Unwrap() string {
	return string(i)
}

type presentImage struct {
	types.Image
}

func (i presentImage) Unwrap() types.Image {
	return i.Image
}

func (i presentImage) Present() bool { return true }

type emptyImage struct{}

func (i emptyImage) Unwrap() types.Image {
	panic("Image not available in empty optional")
}

func (i emptyImage) Present() bool { return false }

func NewPresentImage(i types.Image) types.OptionalImage {
	return presentImage{i}
}

func NewEmptyImage() types.OptionalImage {
	return emptyImage{}
}

type asset struct {
	name  types.AssetName
	image types.OptionalImage
}

func NewAssetWithImage(name types.AssetName, image types.Image) types.Asset {
	return asset{
		name:  name,
		image: presentImage{image},
	}
}

func NewAssetWithoutImage(name types.AssetName) types.Asset {
	return asset{
		name:  name,
		image: emptyImage{},
	}
}

func (a asset) Name() types.AssetName {
	return a.name
}

func (a asset) Image() types.OptionalImage {
	return a.image
}
