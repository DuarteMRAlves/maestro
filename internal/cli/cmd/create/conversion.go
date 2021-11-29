package create

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
)

func MarshalAssetResource(dst *pb.Asset, src *Resource) error {
	if ok, err := assert.ArgStatus(isAssetKind(src), "src not an asset"); !ok {
		return err
	}
	name, nameOk := src.Spec["name"]
	if nameOk {
		dst.Name = name
	}
	image, imageOk := src.Spec["image"]
	if imageOk {
		dst.Image = image
	}
	return nil
}
