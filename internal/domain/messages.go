package domain

import "github.com/DuarteMRAlves/maestro/internal/types"

type CreateAssetRequest struct {
	Name  string
	Image types.OptionalString
}

type CreateAssetResponse struct {
	Err types.OptionalError
}
