package internal

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var assetNameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type AssetName struct{ val string }

var emptyAssetName = AssetName{}

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

type Image struct{ val string }

var emptyImage = Image{val: ""}

func (i Image) Unwrap() string {
	return i.val
}

func (i Image) IsEmpty() bool {
	return i.val == emptyImage.val
}

func NewImage(img string) Image { return Image{val: img} }

type OptionalImage struct {
	val     Image
	present bool
}

func (o OptionalImage) Unwrap() Image {
	return o.val
}

func (o OptionalImage) Present() bool { return o.present }

func NewPresentImage(i Image) OptionalImage {
	return OptionalImage{val: i, present: true}
}

func NewEmptyImage() OptionalImage {
	return OptionalImage{}
}

type Asset struct {
	name  AssetName
	image OptionalImage
}

func (a Asset) Name() AssetName {
	return a.name
}

func (a Asset) Image() OptionalImage {
	return a.image
}

func NewAsset(name AssetName, image OptionalImage) Asset {
	return Asset{
		name:  name,
		image: image,
	}
}
