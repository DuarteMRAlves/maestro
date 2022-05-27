package execute

import "fmt"

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

type logger struct{ debug bool }

func (l logger) Debugf(format string, args ...any) {
	if l.debug {
		fmt.Printf(format, args...)
	}
}

func (l logger) Infof(format string, args ...any) {
	fmt.Printf(format, args...)
}
