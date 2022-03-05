package link

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/types"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStoreWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		link     types.Link
		expected []byte
	}{
		{
			name: "required fields",
			link: &link{
				name: linkName("some-name"),
				source: &linkEndpoint{
					stage: createStage(t, "source", true),
					field: emptyMessageField{},
				},
				target: &linkEndpoint{
					stage: createStage(t, "target", true),
					field: emptyMessageField{},
				},
			},
			expected: []byte("source;;target;"),
		},
		{
			name: "all fields",
			link: &link{
				name: linkName("some-name"),
				source: &linkEndpoint{
					stage: createStage(t, "source", false),
					field: presentMessageField{messageField("source-field")},
				},
				target: &linkEndpoint{
					stage: createStage(t, "target", false),
					field: presentMessageField{messageField("target-field")},
				},
			},
			expected: []byte("source;source-field;target;target-field"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var (
					storedLink   []byte
					loadedSource types.Stage
					loadedTarget types.Stage
				)

				db := kv.NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						store := StoreWithTxn(txn)
						result := store(test.link)
						return result.Error()
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						linkItem, err := txn.Get(kvKey(test.link.Name()))
						if err != nil {
							return err
						}
						storedLink, err = linkItem.ValueCopy(nil)
						if err != nil {
							return err
						}
						loadStage := stage.LoadWithTxn(txn)
						sourceRes := loadStage(test.link.Source().Stage().Name())
						if sourceRes.IsError() {
							return sourceRes.Error()
						}
						loadedSource = sourceRes.Unwrap()
						targetRes := loadStage(test.link.Target().Stage().Name())
						if targetRes.IsError() {
							return targetRes.Error()
						}
						loadedTarget = targetRes.Unwrap()
						return nil
					},
				)
				assert.Equal(t, len(test.expected), len(storedLink))
				for i, e := range test.expected {
					assert.Equal(t, e, storedLink[i])
				}
				assertEqualStage(t, test.link.Source().Stage(), loadedSource)
				assertEqualStage(t, test.link.Target().Stage(), loadedTarget)
			},
		)
	}
}

func TestLoadWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		expected types.Link
		stored   []byte
	}{
		{
			name: "required fields",
			expected: &link{
				name: linkName("some-name"),
				source: &linkEndpoint{
					stage: createStage(t, "source", true),
					field: emptyMessageField{},
				},
				target: &linkEndpoint{
					stage: createStage(t, "target", true),
					field: emptyMessageField{},
				},
			},
			stored: []byte("source;;target;"),
		},
		{
			name: "all fields",
			expected: &link{
				name: linkName("some-name"),
				source: &linkEndpoint{
					stage: createStage(t, "source", false),
					field: presentMessageField{messageField("source-field")},
				},
				target: &linkEndpoint{
					stage: createStage(t, "target", false),
					field: presentMessageField{messageField("target-field")},
				},
			},
			stored: []byte("source;source-field;target;target-field"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded types.Link

				db := kv.NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						storeStage := stage.StoreWithTxn(txn)
						res := storeStage(test.expected.Source().Stage())
						if res.IsError() {
							return res.Error()
						}
						res = storeStage(test.expected.Target().Stage())
						if res.IsError() {
							return res.Error()
						}
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
				fmt.Println(loaded)
				assert.Equal(
					t,
					test.expected.Name().Unwrap(),
					loaded.Name().Unwrap(),
				)
				assertEqualStage(
					t,
					test.expected.Source().Stage(),
					loaded.Source().Stage(),
				)
				assertEqualStage(
					t,
					test.expected.Target().Stage(),
					loaded.Target().Stage(),
				)
				if test.expected.Source().Field().Present() {
					assert.Equal(
						t,
						test.expected.Source().Field().Unwrap(),
						loaded.Source().Field().Unwrap(),
					)
				} else {
					assert.Assert(t, !loaded.Source().Field().Present())
				}
				if test.expected.Target().Field().Present() {
					assert.Equal(
						t,
						test.expected.Target().Field().Unwrap(),
						loaded.Target().Field().Unwrap(),
					)
				} else {
					assert.Assert(t, !loaded.Target().Field().Present())
				}
			},
		)
	}
}

func createStage(t *testing.T, name string, requiredOnly bool) types.Stage {
	var (
		serviceOpt = stage.NewEmptyService()
		methodOpt  = stage.NewEmptyMethod()
	)
	stageName, err := stage.NewStageName(name)
	assert.NilError(t, err, "create name for stage %s", name)

	addr, err := stage.NewAddress(fmt.Sprintf("address-%s", name))
	assert.NilError(t, err, "create address for stage %s", name)

	if !requiredOnly {
		service, err := stage.NewService(fmt.Sprintf("service-%s", name))
		assert.NilError(t, err, "create service for stage %s", name)
		serviceOpt = stage.NewPresentService(service)
		method, err := stage.NewMethod(fmt.Sprintf("method-%s", name))
		assert.NilError(t, err, "create method for stage %s", name)
		methodOpt = stage.NewPresentMethod(method)
	}

	ctx := stage.NewMethodContext(addr, serviceOpt, methodOpt)
	return stage.NewStage(stageName, ctx)
}

func assertEqualStage(t *testing.T, expected types.Stage, actual types.Stage) {
	assert.Equal(
		t,
		expected.Name().Unwrap(),
		actual.Name().Unwrap(),
	)
	assert.Equal(
		t,
		expected.MethodContext().Address().Unwrap(),
		actual.MethodContext().Address().Unwrap(),
	)
	if expected.MethodContext().Service().Present() {
		assert.Equal(
			t,
			expected.MethodContext().Service().Unwrap().Unwrap(),
			actual.MethodContext().Service().Unwrap().Unwrap(),
		)
	}
	if expected.MethodContext().Method().Present() {
		assert.Equal(
			t,
			expected.MethodContext().Method().Unwrap().Unwrap(),
			actual.MethodContext().Method().Unwrap().Unwrap(),
		)
	}
}
