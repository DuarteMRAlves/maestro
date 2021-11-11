package asset

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

const IdSize = 10

// Asset represents components that can be orchestrated
// They have a protobuf definition and an associated
// docker image to be executed
type Asset struct {
	Id    identifier.Id
	Name  string
	Image string
}

func (a *Asset) Clone() *Asset {
	return &Asset{
		Id:    a.Id,
		Name:  a.Name,
		Image: a.Image,
	}
}

func (a *Asset) String() string {
	return fmt.Sprintf("Asset{Id:%v,Name:%v,Image:%v}", a.Id, a.Name, a.Image)
}
