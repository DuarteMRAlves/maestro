package kv

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func NewDb() (*badger.DB, error) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, errdefs.InternalWithMsg("create database: %v", err)
	}
	return db, err
}

func NewTestDb(t *testing.T) *badger.DB {
	opts := badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	assert.NilError(t, err, "error creating test db")
	return db
}
