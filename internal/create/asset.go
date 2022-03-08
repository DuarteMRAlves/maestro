package create

import "github.com/DuarteMRAlves/maestro/internal/domain"

type AssetSaver interface {
	Save(domain.Asset) domain.AssetResult
}

type AssetLoader interface {
	Load(domain.AssetName) domain.AssetResult
}

type AssetExistsVerifier interface {
	Verify(domain.AssetName) bool
}

type AssetStorage interface {
	AssetSaver
	AssetLoader
	AssetExistsVerifier
}

type AssetRequest struct {
	Name  string
	Image domain.OptionalString
}

type AssetResponse struct {
	Err domain.OptionalError
}
