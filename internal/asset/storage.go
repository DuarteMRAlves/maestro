package asset

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"strings"
)

func StoreAssetWithTxn(txn *badger.Txn) func(domain.Asset) Result {
	return func(a domain.Asset) Result {
		var (
			buf bytes.Buffer
			err error
		)
		buf.WriteString(a.Name().Unwrap())
		buf.WriteByte(';')
		buf.WriteString(a.Image().Unwrap())
		err = txn.Set(kv.AssetKey(a.Name()), buf.Bytes())
		if err != nil {
			err = errdefs.InternalWithMsg("storage error: %v", err)
			return NewErrResult(err)
		}
		return NewResult(a)
	}
}

func LoadAssetWithTxn(txn *badger.Txn) func(domain.AssetName) Result {
	return func(name domain.AssetName) Result {
		var (
			item *badger.Item
			data []byte
			err  error
		)
		item, err = txn.Get(kv.AssetKey(name))
		if err != nil {
			err = errdefs.PrependMsg(err, "load asset %s", name)
			return NewErrResult(err)
		}

		data, err = item.ValueCopy(nil)
		buf := bytes.NewBuffer(data)
		splits := strings.Split(buf.String(), ";")
		if len(splits) != 2 {
			err := errdefs.InternalWithMsg(
				"invalid format: expected 2 semi-colon separated values",
			)
			return NewErrResult(err)
		}
		name, err = NewAssetName(splits[0])
		if err != nil {
			err = errdefs.PrependMsg(err, "load asset %s", name)
			return NewErrResult(err)
		}
		img := NewImage(splits[1])
		a := NewAsset(name, img)
		return NewResult(a)
	}
}
