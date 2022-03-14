package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"gotest.tools/v3/assert"
	"testing"
)

func NewDb() (*badger.DB, error) {
	lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
	zapLogger, err := logs.DefaultProductionLogger(lvl)
	if err != nil {
		return nil, err
	}
	logger := NewBadgerLogger(zapLogger.Sugar())
	opts := badger.DefaultOptions("").WithInMemory(true).WithLogger(logger)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, err
}

func NewTestDb(t *testing.T) *badger.DB {
	lvl := zap.NewAtomicLevelAt(zap.WarnLevel)
	zapLogger, err := logs.DefaultProductionLogger(lvl)
	assert.NilError(t, err, "error creating db logger")
	logger := NewBadgerLogger(zapLogger.Sugar())
	opts := badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.WARNING).
		WithLogger(logger)
	db, err := badger.Open(opts)
	assert.NilError(t, err, "error creating test db")
	return db
}

type badgerLogger struct {
	logger *zap.SugaredLogger
}

func NewBadgerLogger(logger *zap.SugaredLogger) badger.Logger {
	return &badgerLogger{logger: logger}
}

func (b *badgerLogger) Errorf(s string, i ...interface{}) {
	b.logger.Errorf(s, i...)
}

func (b *badgerLogger) Warningf(s string, i ...interface{}) {
	b.logger.Warnf(s, i...)
}

func (b *badgerLogger) Infof(s string, i ...interface{}) {
	b.logger.Infof(s, i...)
}

func (b *badgerLogger) Debugf(s string, i ...interface{}) {
	b.logger.Debugf(s, i...)
}
