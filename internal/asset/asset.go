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

func NewImage(img string) domain.Image {
	return image(img)
}

func (i image) Unwrap() string {
	return string(i)
}

type asset struct {
	name  domain.AssetName
	image domain.Image
}

func NewAsset(name domain.AssetName, image domain.Image) domain.Asset {
	return asset{
		name:  name,
		image: image,
	}
}

func (a asset) Name() domain.AssetName {
	return a.name
}

func (a asset) Image() domain.Image {
	return a.image
}
