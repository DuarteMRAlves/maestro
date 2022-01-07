package flow

import "github.com/DuarteMRAlves/maestro/internal/link"

// Input represents the several input flows for a stage
type Input struct {
	connections map[string]*link.Link
}

func NewInput() *Input {
	return &Input{connections: map[string]*link.Link{}}
}
