package stage

import "fmt"

type AlreadyExists struct {
	Name string
}

func (e AlreadyExists) Error() string {
	return fmt.Sprintf("Stage '%v' already exists.", e.Name)
}
