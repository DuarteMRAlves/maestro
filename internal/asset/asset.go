package asset

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

// Asset represents components that can be orchestrated
// They have a protobuf definition and an associated
// docker image to be executed
type Asset struct {
	name  apitypes.AssetName
	image string
}

func New(name apitypes.AssetName, image string) *Asset {
	return &Asset{
		name:  name,
		image: image,
	}
}

func (a *Asset) Name() apitypes.AssetName {
	return a.name
}

func (a *Asset) Image() string {
	return a.image
}

func (a *Asset) Clone() *Asset {
	return &Asset{
		name:  a.name,
		image: a.image,
	}
}

func (a *Asset) ToApi() *apitypes.Asset {
	return &apitypes.Asset{
		Name:  a.name,
		Image: a.image,
	}
}

func (a Asset) String() string {
	return fmt.Sprintf("Asset{name:%v,image:%v}", a.name, a.image)
}
