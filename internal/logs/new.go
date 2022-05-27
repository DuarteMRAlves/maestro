package logs

import (
	"io"
	"log"
)

// Logger provides a simple wrapper around log.Logger to print messages at two
// levels: info or debug. info displays information for the user and debug for
// displays information for the developer.
type Logger struct {
	logger *log.Logger
	debug  bool
}

// New creates a logger wrapping the default log.Logger. debug specifies whether
// debug messages should be output or not.
func New(debug bool) Logger {
	return NewWithLogger(defaultLogger(), debug)
}

// NewWithOutput is equivalent to New with a custom out writer.
func NewWithOutput(w io.Writer, debug bool) Logger {
	return NewWithLogger(log.New(w, "", log.LstdFlags), debug)
}

// NewWithLogger is equivalent to New with a custom logger. If logger is nil
// the default will be used.
func NewWithLogger(logger *log.Logger, debug bool) Logger {
	return Logger{
		logger: logger,
		debug:  debug,
	}
}

func (l Logger) Debugf(format string, args ...any) {
	if l.debug {
		l.writef(format, args...)
	}
}

func (l Logger) Infof(format string, args ...any) {
	l.writef(format, args...)
}

func (l Logger) writef(format string, args ...any) {
	logger := l.logger
	if logger == nil {
		logger = defaultLogger()
	}
	logger.Printf(format, args...)
}

func defaultLogger() *log.Logger {
	return log.Default()
}
