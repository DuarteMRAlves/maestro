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

func MarshalStageResource(dst *pb.Stage, src *Resource) error {
	if ok, err := assert.ArgStatus(isStageKind(src), "src not an stage"); !ok {
		return err
	}
	name, nameOk := src.Spec["name"]
	if nameOk {
		dst.Name = name
	}
	asset, assetOk := src.Spec["asset"]
	if assetOk {
		dst.Asset = asset
	}
	service, serviceOk := src.Spec["service"]
	if serviceOk {
		dst.Service = service
	}
	method, methodOk := src.Spec["method"]
	if methodOk {
		dst.Method = method
	}
	return nil
}
