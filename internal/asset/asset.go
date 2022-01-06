package asset

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

const IdSize = 10

// Asset represents components that can be orchestrated
// They have a protobuf definition and an associated
// docker image to be executed
type Asset struct {
	Name  string
	Image string
}

func (a *Asset) Clone() *Asset {
	return &Asset{
		Name:  a.Name,
		Image: a.Image,
	}
}

func (a *Asset) ToApi() *apitypes.Asset {
	return &apitypes.Asset{
		Name:  a.Name,
		Image: a.Image,
	}
}

func (a *Asset) String() string {
	return fmt.Sprintf("Asset{Name:%v,Image:%v}", a.Name, a.Image)
}
