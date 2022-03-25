package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
	"testing"
)

func NewDb() (*badger.DB, error) {
	logger := NewBadgerLogger(logs.New(false))
	opts := badger.DefaultOptions("").WithInMemory(true).WithLogger(logger)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, err
}

func NewTestDb(t *testing.T) *badger.DB {
	logger := NewBadgerLogger(logs.New(true))
	opts := badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.WARNING).
		WithLogger(logger)
	db, err := badger.Open(opts)
	if err == nil {
		t.Fatalf("error creating test db: %s", err)
	}
	return db
}

type badgerLogger struct {
	logger logs.Logger
}

func NewBadgerLogger(logger logs.Logger) badger.Logger {
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
