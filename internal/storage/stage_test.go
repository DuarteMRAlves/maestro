package storage

//
// import (
// 	"github.com/DuarteMRAlves/maestro/internal/create"
// 	"github.com/dgraph-io/badger/v3"
// 	"gotest.tools/v3/assert"
// 	"testing"
// )
//
// func TestSaveStageWithTxn(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		stage    create.Stage
// 		expected []byte
// 	}{
// 		{
// 			name:     "required fields",
// 			stage:    createStage(t, "some-name", true),
// 			expected: []byte("address-some-name;;"),
// 		},
// 		{
// 			name:     "all fields",
// 			stage:    createStage(t, "some-name", false),
// 			expected: []byte("address-some-name;service-some-name;method-some-name"),
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(
// 			test.name,
// 			func(t *testing.T) {
// 				var stored []byte
//
// 				db := NewTestDb(t)
// 				defer db.Close()
//
// 				err := db.Update(
// 					func(txn *badger.Txn) error {
// 						store := SaveStageWithTxn(txn)
// 						result := store(test.stage)
// 						return result.Error()
// 					},
// 				)
// 				assert.NilError(t, err, "save error")
//
// 				err = db.View(
// 					func(txn *badger.Txn) error {
// 						item, err := txn.Get(kvKey(test.stage.Name()))
// 						if err != nil {
// 							return err
// 						}
// 						stored, err = item.ValueCopy(nil)
// 						return err
// 					},
// 				)
// 				assert.Equal(t, len(test.expected), len(stored), "stored size")
// 				for i, e := range test.expected {
// 					assert.Equal(t, e, stored[i], "stored not equal")
// 				}
// 			},
// 		)
// 	}
// }
//
// func TestLoadStageWithTxn(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected create.Stage
// 		stored   []byte
// 	}{
// 		{
// 			name:     "required fields",
// 			expected: createStage(t, "some-name", true),
// 			stored:   []byte("address-some-name;;"),
// 		},
// 		{
// 			name:     "all fields",
// 			expected: createStage(t, "some-name", false),
// 			stored:   []byte("address-some-name;service-some-name;method-some-name"),
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(
// 			test.name,
// 			func(t *testing.T) {
// 				var loaded create.Stage
//
// 				db := NewTestDb(t)
// 				defer db.Close()
//
// 				err := db.Update(
// 					func(txn *badger.Txn) error {
// 						return txn.Set(kvKey(test.expected.Name()), test.stored)
// 					},
// 				)
// 				assert.NilError(t, err, "save error")
//
// 				err = db.View(
// 					func(txn *badger.Txn) error {
// 						load := LoadStageWithTxn(txn)
// 						res := load(test.expected.Name())
// 						if !res.IsError() {
// 							loaded = res.Unwrap()
// 						}
// 						return res.Error()
// 					},
// 				)
// 				assert.NilError(t, err, "load error")
// 				assert.Equal(
// 					t,
// 					test.expected.Name().Unwrap(),
// 					loaded.Name().Unwrap(),
// 				)
// 				assert.Equal(
// 					t,
// 					test.expected.MethodContext().Address().Unwrap(),
// 					loaded.MethodContext().Address().Unwrap(),
// 				)
// 				if test.expected.MethodContext().Service().Present() {
// 					assert.Equal(
// 						t,
// 						test.expected.MethodContext().Service().Unwrap().Unwrap(),
// 						loaded.MethodContext().Service().Unwrap().Unwrap(),
// 					)
// 				}
// 				if test.expected.MethodContext().Method().Present() {
// 					assert.Equal(
// 						t,
// 						test.expected.MethodContext().Method().Unwrap().Unwrap(),
// 						loaded.MethodContext().Method().Unwrap().Unwrap(),
// 					)
// 				}
// 			},
// 		)
// 	}
// }
