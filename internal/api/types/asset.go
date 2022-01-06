package types

// Asset represents an image with a grpc server that can be deployed.
type Asset struct {
	// Name that should be associated with the asset. Is required and should be
	// unique.
	// (required, unique)
	Name AssetName `yaml:"name" info:"required"`
	// Image specifies the container image that should be associated with this
	// asset
	// (optional)
	Image string `yaml:"image"`
}

// AssetName is a name that uniquely identifies an Asset
type AssetName string
