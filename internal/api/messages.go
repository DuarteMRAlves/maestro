package api

import apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"

// CreateAssetRequest represents a message to create an asset.
type CreateAssetRequest struct {
	// Name that should be associated with the asset. Is required and should be
	// unique.
	// (required, unique)
	Name apitypes.AssetName `yaml:"name" info:"required"`
	// Image specifies the container image that should be associated with this
	// asset
	// (optional)
	Image string `yaml:"image"`
}

// GetAssetRequest represents a message to get assets with specific
// characteristics.
type GetAssetRequest struct {
	// Name should be set to retrieve only assets with the given name.
	Name apitypes.AssetName
	// Image should be set to retrieve only assets with the given image.
	Image string
}
