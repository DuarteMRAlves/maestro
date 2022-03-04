package asset

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/types"
	"github.com/dgraph-io/badger/v3"
	"strings"
)

func StoreAssetWithTxn(txn *badger.Txn) func(types.Asset) types.AssetResult {
	return func(a types.Asset) types.AssetResult {
		var (
			buf bytes.Buffer
			err error
		)
		buf.WriteString(a.Name().Unwrap())
		buf.WriteByte(';')
		buf.WriteString(imageToString(a.Image()))
		err = txn.Set(kvKey(a.Name()), buf.Bytes())
		if err != nil {
			err = errdefs.InternalWithMsg("storage error: %v", err)
			return NewErrResult(err)
		}
		return NewResult(a)
	}
}

func LoadAssetWithTxn(txn *badger.Txn) func(types.AssetName) types.AssetResult {
	return func(name types.AssetName) types.AssetResult {
		var (
			item *badger.Item
			data []byte
			err  error
		)
		item, err = txn.Get(kvKey(name))
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

func imageToString(img types.OptionalImage) string {
	if img.Present() {
		return img.Unwrap().Unwrap()
	} else {
		return ""
	}
}

func stringToImage(data string) (types.OptionalImage, error) {
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

func kvKey(name types.AssetName) []byte {
	return []byte(fmt.Sprintf("asset:%s", name.Unwrap()))
}
