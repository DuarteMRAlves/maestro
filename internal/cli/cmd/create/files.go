package create

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"time"
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	if createAll || kind == resources.AssetKind {
		for _, r := range parsed {
			if resources.IsAssetKind(r) {
				a := &resources.AssetResource{}
				if err = resources.MarshalResource(a, r); err != nil {
					return err
				}
				if err = client.CreateAsset(ctx, a, addr); err != nil {
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
				if err = client.CreateStage(ctx, s, addr); err != nil {
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
				if err = client.CreateLink(ctx, l, addr); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
