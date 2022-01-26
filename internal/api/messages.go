package api

// CreateOrchestrationRequest represents a message to create an orchestration.
type CreateOrchestrationRequest struct {
	// Name that should be associated with the orchestration.
	// (required, unique)
	Name OrchestrationName `yaml:"name" info:"required"`
	// Links specifies the name of the links that compose this orchestration.
	// (required, non-empty)
	Links []LinkName `yaml:"links" info:"required"`
}

// GetOrchestrationRequest represents a message to retrieve orchestrations with
// specific characteristics.
type GetOrchestrationRequest struct {
	// Name should be set to retrieve orchestrations with the given name.
	Name OrchestrationName
	// Phase should be set to retrieve orchestrations in a given phase.
	Phase OrchestrationPhase
}

// CreateAssetRequest represents a message to create an asset.
type CreateAssetRequest struct {
	// Name that should be associated with the asset. Is required and should be
	// unique.
	// (required, unique)
	Name AssetName `yaml:"name" info:"required"`
	// Image specifies the container image that should be associated with this
	// asset
	// (optional)
	Image string `yaml:"image"`
}

// GetAssetRequest represents a message to get assets with specific
// characteristics.
type GetAssetRequest struct {
	// Name should be set to retrieve only assets with the given name.
	Name AssetName
	// Image should be set to retrieve only assets with the given image.
	Image string
}
