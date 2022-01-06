package output

import "github.com/DuarteMRAlves/maestro/internal/link"

// Output represents the several input fields
type Output struct {
	connections map[string]*link.Link
}

func NewDefault() *Output {
	return &Output{connections: map[string]*link.Link{}}
}
