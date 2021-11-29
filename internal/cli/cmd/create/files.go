package create

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
)

func createFromFiles(files []string, addr string, kind string) error {
	resources, err := ParseResources(files)
	if err != nil {
		return err
	}

	if err = isValidKinds(resources); err != nil {
		return err
	}

	createAll := kind == ""

	if createAll || kind == assetKind {
		for _, r := range resources {
			if isAssetKind(r) {
				a := &pb.Asset{}
				if err = MarshalAssetResource(a, r); err != nil {
					return err
				}
				if err = client.CreateAsset(a, addr); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
