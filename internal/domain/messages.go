package domain

type CreateAssetRequest struct {
	Name  string
	Image OptionalString
}

type CreateAssetResponse struct {
	Err OptionalError
}
