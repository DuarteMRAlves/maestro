package storage

import "github.com/dgraph-io/badger/v3"

func NewDb(logger Logger) (*badger.DB, error) {
	badgerLogger := NewBadgerLogger(logger)
	opts := badger.DefaultOptions("").WithInMemory(true).WithLogger(badgerLogger)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, err
}

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

type badgerLogger struct {
	logger Logger
}

func NewBadgerLogger(logger Logger) badger.Logger {
	return &badgerLogger{logger: logger}
}

func (b *badgerLogger) Errorf(s string, i ...interface{}) {
	b.logger.Infof(s, i...)
}

func (b *badgerLogger) Warningf(s string, i ...interface{}) {
	b.logger.Infof(s, i...)
}

func (b *badgerLogger) Infof(s string, i ...interface{}) {
	b.logger.Infof(s, i...)
}

func (b *badgerLogger) Debugf(s string, i ...interface{}) {
	b.logger.Debugf(s, i...)
}
