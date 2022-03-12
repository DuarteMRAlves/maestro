package internal

import "fmt"

type AssetNotFound struct {
	Name AssetName
}

func (err *AssetNotFound) Error() string {
	return fmt.Sprintf("asset '%s' not found", err.Name)
}
