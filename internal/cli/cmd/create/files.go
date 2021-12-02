package create

import (
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	resources "github.com/DuarteMRAlves/maestro/internal/cli/resources"
)

func createFromFiles(files []string, addr string, kind string) error {
	parsed, err := resources.ParseFiles(files)
	if err != nil {
		return err
	}
	if err = resources.IsValidKinds(parsed); err != nil {
		return err
	}

	createAll := kind == ""

	if createAll || kind == resources.AssetKind {
		for _, r := range parsed {
			if resources.IsAssetKind(r) {
				a := &resources.AssetResource{}
				if err = resources.MarshalResource(a, r); err != nil {
					return err
				}
				if err = client.CreateAsset(a, addr); err != nil {
					return err
				}
			}
		}
	}
	if createAll || kind == resources.StageKind {
		for _, r := range parsed {
			if resources.IsStageKind(r) {
				s := &resources.StageResource{}
				if err = resources.MarshalResource(s, r); err != nil {
					return err
				}
				if err = client.CreateStage(s, addr); err != nil {
					return err
				}
			}
		}
	}
	if createAll || kind == resources.LinkKind {
		for _, r := range parsed {
			if resources.IsLinkKind(r) {
				l := &resources.LinkResource{}
				if err = resources.MarshalResource(l, r); err != nil {
					return err
				}
				if err = client.CreateLink(l, addr); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
