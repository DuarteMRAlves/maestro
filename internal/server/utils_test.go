package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/util"
)

// assetForNum deterministically creates an asset with the given number.
func assetForNum(num int) *api.Asset {
	name := util.AssetNameForNum(num)
	img := util.AssetImageForNum(num)
	return &api.Asset{
		Name:  name,
		Image: img,
	}
}
