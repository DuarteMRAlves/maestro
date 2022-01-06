package input

import "github.com/DuarteMRAlves/maestro/internal/link"

// Input represents the several input fields
type Input struct {
	connections map[string]*link.Link
}

func NewDefault() *Input {
	return &Input{connections: map[string]*link.Link{}}
}
