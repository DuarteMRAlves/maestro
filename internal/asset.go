package internal

import (
	"fmt"
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

func (a AssetName) String() string {
	return a.val
}

func NewAssetName(name string) (AssetName, error) {
	if isValidAssetName(name) {
		return AssetName{name}, nil
	}
	return emptyAssetName, &invalidAssetName{name: name}
}

func isValidAssetName(name string) bool {
	return assetNameRegExp.MatchString(name)
}

type invalidAssetName struct{ name string }

func (err *invalidAssetName) Error() string {
	return fmt.Sprintf("invalid asset name: '%s'", err.name)
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

type Asset struct {
	name  AssetName
	image Image
}

func (a Asset) Name() AssetName {
	return a.name
}

func (a Asset) Image() Image {
	return a.image
}

func NewAsset(name AssetName, image Image) Asset {
	return Asset{
		name:  name,
		image: image,
	}
}
