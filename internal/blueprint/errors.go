package blueprint

import "fmt"

type AlreadyExists struct {
	Name string
}

func (e AlreadyExists) Error() string {
	return fmt.Sprintf("Blueprint with Name '%v' already exists", e.Name)
}
