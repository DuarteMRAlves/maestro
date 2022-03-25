package mock

import "fmt"

type Logger struct {
	DebugActive bool
	InfoActive  bool
}

func (l *Logger) Debugf(format string, args ...any) {
	if l.DebugActive {
		format = fmt.Sprintf("debug\t%s", format)
		fmt.Printf(format, args...)
	}
}

func (l *Logger) Infof(format string, args ...any) {
	if l.InfoActive || l.DebugActive {
		format = fmt.Sprintf("info\t%s", format)
		fmt.Printf(format, args...)
	}
}
