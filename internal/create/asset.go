package create

import "github.com/DuarteMRAlves/maestro/internal/domain"

type SaveAsset func(domain.Asset) domain.AssetResult
type LoadAsset func(domain.AssetName) domain.AssetResult
type ExistsAsset func(domain.AssetName) bool

type AssetRequest struct {
	Name  string
	Image domain.OptionalString
}

type AssetResponse struct {
	Err domain.OptionalError
}
