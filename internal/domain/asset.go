package domain

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var assetNameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type AssetName struct{ val string }

var emptyAssetName = AssetName{val: ""}

func (a AssetName) Unwrap() string {
	return a.val
}

func (a AssetName) IsEmpty() bool {
	return a.val == emptyAssetName.val
}

func NewAssetName(name string) (AssetName, error) {
	if isValidAssetName(name) {
		return AssetName{name}, nil
	}
	err := errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
	return emptyAssetName, err
}

func isValidAssetName(name string) bool {
	return assetNameRegExp.MatchString(name)
}

type image string

func (i image) Image() {}

func (i image) Unwrap() string {
	return string(i)
}

func NewImage(img string) (Image, error) {
	if len(img) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty image")
	}
	return image(img), nil
}

type presentImage struct {
	Image
}

func (i presentImage) Unwrap() Image {
	return i.Image
}

func (i presentImage) Present() bool { return true }

type emptyImage struct{}

func (i emptyImage) Unwrap() Image {
	panic("Image not available in empty optional")
}

func (i emptyImage) Present() bool { return false }

func NewPresentImage(i Image) OptionalImage {
	return presentImage{i}
}

func NewEmptyImage() OptionalImage {
	return emptyImage{}
}

type asset struct {
	name  AssetName
	image OptionalImage
}

func (a asset) Name() AssetName {
	return a.name
}

func (a asset) Image() OptionalImage {
	return a.image
}

func NewAsset(name AssetName, image OptionalImage) Asset {
	return asset{
		name:  name,
		image: image,
	}
}
