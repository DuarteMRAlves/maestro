package mock

import "fmt"

type notFound struct {
	typ  string
	name string
}

func (err *notFound) NotFound() { /* Do nothing. Just implement NotFound interfaces */ }

func (err *notFound) Error() string {
	return fmt.Sprintf("%s '%s' not found.", err.typ, err.name)
}
