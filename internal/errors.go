package internal

import "fmt"

// NotFound defines an error when a resource does not exist.
type NotFound struct {
	// Type defines the resource type (e.g.: file, orchestration, image, etc.).
	Type string
	// Ident is the concrete resource identifier that was requested. (e.g.
	// file.tsv, orchestration-1, etc.).
	Ident string
}

func (err *NotFound) Error() string {
	return fmt.Sprintf("%s '%s' not found.", err.Type, err.Ident)
}

// AlreadyExists defines an error when a resource already exists.
type AlreadyExists struct {
	// Type defines the resource type (e.g.: file, orchestration, image, etc.).
	Type string
	// Ident is the concrete resource identifier that was requested. (e.g.
	// file.tsv, orchestration-1, etc.).
	Ident string
}

func (err *AlreadyExists) Error() string {
	return fmt.Sprintf("%s '%s' already exists.", err.Type, err.Ident)
}
