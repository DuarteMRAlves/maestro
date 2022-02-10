package storage

import (
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func NewTestDb(t *testing.T) *badger.DB {
	opts := badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	assert.NilError(t, err, "error creating test db")
	return db
}
