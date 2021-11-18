package asset

import "fmt"

type AlreadyExists struct {
	Name string
}

func (e AlreadyExists) Error() string {
	return fmt.Sprintf("Asset with Name '%v' already exists", e.Name)
}
