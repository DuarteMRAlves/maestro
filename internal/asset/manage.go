package asset

import (
	"bytes"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/dgraph-io/badger/v3"
)

func Create(txn *badger.Txn, cfg *apitypes.Asset) error {
	var err error
	if err = validateCreateAssetConfig(cfg); err != nil {
		return errdefs.InvalidArgumentWithError(err)
	}

	if Contains(txn, cfg.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"asset '%v' already exists",
			cfg.Name,
		)
	}
	asset := New(cfg.Name, cfg.Image)
	if err = Persist(txn, asset); err != nil {
		return errdefs.InternalWithMsg("persist error: %v", err)
	}
	return nil
}

func Contains(txn *badger.Txn, name apitypes.AssetName) bool {
	item, _ := txn.Get(assetKey(name))
	return item != nil
}

func Get(txn *badger.Txn, query *apitypes.Asset) ([]*apitypes.Asset, error) {
	var (
		asset Asset
		cp    []byte
		err   error
	)

	if query == nil {
		query = &apitypes.Asset{}
	}
	filter := buildQueryFilter(query)
	res := make([]*apitypes.Asset, 0)
	it := txn.NewIterator(badger.DefaultIteratorOptions)

	defer it.Close()
	prefix := []byte("asset:")

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		cp, err = item.ValueCopy(cp)
		if err != nil {
			return nil, errdefs.InternalWithMsg("read: %v", err)
		}
		err = load(&asset, cp)
		if err != nil {
			return nil, errdefs.InternalWithMsg("decoding: %v", err)
		}
		if filter(&asset) {
			res = append(res, asset.ToApi())
		}
	}
	return res, nil
}

func buildQueryFilter(query *apitypes.Asset) func(a *Asset) bool {
	filters := make([]func(a *Asset) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(a *Asset) bool {
				return a.Name() == query.Name
			},
		)
	}
	if query.Image != "" {
		filters = append(
			filters,
			func(a *Asset) bool {
				return a.Image() == query.Image
			},
		)
	}
	if len(filters) > 0 {
		return func(a *Asset) bool {
			for _, f := range filters {
				if !f(a) {
					return false
				}
			}
			return true
		}
	}
	return func(a *Asset) bool {
		return true
	}
}

// validateCreateAssetConfig verifies if all conditions to create an asset are
// met. It returns an error if a condition is not met and nil otherwise.
func validateCreateAssetConfig(cfg *apitypes.Asset) error {
	if ok, err := validate.ArgNotNil(cfg, "cfg"); !ok {
		return errdefs.InvalidArgumentWithError(err)
	}
	if !naming.IsValidAssetName(cfg.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			cfg.Name,
		)
	}
	return nil
}

func Persist(txn *badger.Txn, a *Asset) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(&buf, a.name, a.image)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = txn.Set(assetKey(a.name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func load(a *Asset, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &a.name, &a.image)
	return err
}

// assetKey returns the image key for an asset with the given name
func assetKey(name apitypes.AssetName) []byte {
	return []byte(fmt.Sprintf("asset:%s", name))
}
