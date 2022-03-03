package asset

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"strings"
)

func StoreAssetWithTxn(txn *badger.Txn) func(domain.Asset) domain.AssetResult {
	return func(a domain.Asset) domain.AssetResult {
		var (
			buf bytes.Buffer
			err error
		)
		buf.WriteString(a.Name().Unwrap())
		buf.WriteByte(';')
		buf.WriteString(imageToString(a.Image()))
		err = txn.Set(kv.AssetKey(a.Name()), buf.Bytes())
		if err != nil {
			err = errdefs.InternalWithMsg("storage error: %v", err)
			return NewErrResult(err)
		}
		return NewResult(a)
	}
}

func LoadAssetWithTxn(txn *badger.Txn) func(domain.AssetName) domain.AssetResult {
	return func(name domain.AssetName) domain.AssetResult {
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
		img, err := stringToImage(splits[1])
		if err != nil {
			err = errdefs.PrependMsg(err, "load asset %s", name)
			return NewErrResult(err)
		}
		if img.Present() {
			return NewResult(NewAssetWithImage(name, img.Unwrap()))
		} else {
			return NewResult(NewAssetWithoutImage(name))
		}
	}
}

func imageToString(img domain.OptionalImage) string {
	if img.Present() {
		return img.Unwrap().Unwrap()
	} else {
		return ""
	}
}

func stringToImage(data string) (domain.OptionalImage, error) {
	if data == "" {
		return NewEmptyImage(), nil
	} else {
		img, err := NewImage(data)
		if err != nil {
			return nil, err
		}
		return NewPresentImage(img), nil
	}
}
