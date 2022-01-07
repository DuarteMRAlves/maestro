package flow

import "github.com/DuarteMRAlves/maestro/internal/link"

// Output represents the several output flows for a stage
type Output struct {
	connections map[string]*link.Link
}

func NewOutput() *Output {
	return &Output{connections: map[string]*link.Link{}}
}
