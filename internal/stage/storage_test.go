package stage

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStoreWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		stage    domain.Stage
		expected []byte
	}{
		{
			name: "required fields",
			stage: stage{
				name: stageName("some-name"),
				methodCtx: methodContext{
					address: address("some-address"),
					service: emptyService{},
					method:  emptyMethod{},
				},
			},
			expected: []byte("some-address;;"),
		},
		{
			name: "all fields",
			stage: stage{
				name: stageName("some-name"),
				methodCtx: methodContext{
					address: address("some-address"),
					service: presentService{service("some-service")},
					method:  presentMethod{method("some-method")},
				},
			},
			expected: []byte("some-address;some-service;some-method"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var stored []byte

				db := kv.NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						store := StoreWithTxn(txn)
						result := store(test.stage)
						return result.Error()
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(kvKey(test.stage.Name()))
						if err != nil {
							return err
						}
						stored, err = item.ValueCopy(nil)
						return err
					},
				)
				assert.Equal(t, len(test.expected), len(stored), "stored size")
				for i, e := range test.expected {
					assert.Equal(t, e, stored[i], "stored not equal")
				}
			},
		)
	}
}

func TestLoadWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		expected domain.Stage
		stored   []byte
	}{
		{
			name: "required fields",
			expected: stage{
				name: stageName("some-name"),
				methodCtx: methodContext{
					address: address("some-address"),
					service: emptyService{},
					method:  emptyMethod{},
				},
			},
			stored: []byte("some-address;;"),
		},
		{
			name: "all fields",
			expected: stage{
				name: stageName("some-name"),
				methodCtx: methodContext{
					address: address("some-address"),
					service: presentService{service("some-service")},
					method:  presentMethod{method("some-method")},
				},
			},
			stored: []byte("some-address;some-service;some-method"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded domain.Stage

				db := kv.NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(kvKey(test.expected.Name()), test.stored)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						load := LoadWithTxn(txn)
						res := load(test.expected.Name())
						if !res.IsError() {
							loaded = res.Unwrap()
						}
						return res.Error()
					},
				)
				assert.NilError(t, err, "load error")
				assert.Equal(
					t,
					test.expected.Name().Unwrap(),
					loaded.Name().Unwrap(),
				)
				assert.Equal(
					t,
					test.expected.MethodContext().Address().Unwrap(),
					loaded.MethodContext().Address().Unwrap(),
				)
				if test.expected.MethodContext().Service().Present() {
					assert.Equal(
						t,
						test.expected.MethodContext().Service().Unwrap().Unwrap(),
						loaded.MethodContext().Service().Unwrap().Unwrap(),
					)
				}
				if test.expected.MethodContext().Method().Present() {
					assert.Equal(
						t,
						test.expected.MethodContext().Method().Unwrap().Unwrap(),
						loaded.MethodContext().Method().Unwrap().Unwrap(),
					)
				}
			},
		)
	}
}
