package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type assetName string

func NewAssetName(name string) (domain.AssetName, error) {
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

func NewImage(img string) (domain.Image, error) {
	if len(img) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty image")
	}
	return image(img), nil
}

func (i image) Unwrap() string {
	return string(i)
}

type presentImage struct {
	domain.Image
}

func (i presentImage) Unwrap() domain.Image {
	return i.Image
}

func (i presentImage) Present() bool { return true }

type emptyImage struct{}

func (i emptyImage) Unwrap() domain.Image {
	panic("Image not available in empty optional")
}

func (i emptyImage) Present() bool { return false }

func NewPresentImage(i domain.Image) domain.OptionalImage {
	return presentImage{i}
}

func NewEmptyImage() domain.OptionalImage {
	return emptyImage{}
}

type asset struct {
	name  domain.AssetName
	image domain.OptionalImage
}

func NewAssetWithImage(name domain.AssetName, image domain.Image) domain.Asset {
	return asset{
		name:  name,
		image: presentImage{image},
	}
}

func NewAssetWithoutImage(name domain.AssetName) domain.Asset {
	return asset{
		name:  name,
		image: emptyImage{},
	}
}

func (a asset) Name() domain.AssetName {
	return a.name
}

func (a asset) Image() domain.OptionalImage {
	return a.image
}
