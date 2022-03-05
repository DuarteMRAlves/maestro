package asset

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
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
		err = txn.Set(kvKey(a.Name()), buf.Bytes())
		if err != nil {
			err = errdefs.InternalWithMsg("storage error: %v", err)
			return domain.ErrAsset(err)
		}
		return domain.SomeAsset(a)
	}
}

func LoadAssetWithTxn(txn *badger.Txn) func(domain.AssetName) domain.AssetResult {
	return func(name domain.AssetName) domain.AssetResult {
		var (
			item *badger.Item
			data []byte
			err  error
		)
		item, err = txn.Get(kvKey(name))
		if err != nil {
			err = errdefs.PrependMsg(err, "load asset %s", name)
			return domain.ErrAsset(err)
		}

		data, err = item.ValueCopy(nil)
		buf := bytes.NewBuffer(data)
		splits := strings.Split(buf.String(), ";")
		if len(splits) != 2 {
			err := errdefs.InternalWithMsg(
				"invalid format: expected 2 semi-colon separated values",
			)
			return domain.ErrAsset(err)
		}
		name, err = domain.NewAssetName(splits[0])
		if err != nil {
			err = errdefs.PrependMsg(err, "load asset %s", name)
			return domain.ErrAsset(err)
		}
		img, err := stringToImage(splits[1])
		if err != nil {
			err = errdefs.PrependMsg(err, "load asset %s", name)
			return domain.ErrAsset(err)
		}
		if img.Present() {
			return domain.SomeAsset(
				domain.NewAssetWithImage(
					name,
					img.Unwrap(),
				),
			)
		} else {
			return domain.SomeAsset(domain.NewAssetWithoutImage(name))
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
		return domain.NewEmptyImage(), nil
	} else {
		img, err := domain.NewImage(data)
		if err != nil {
			return nil, err
		}
		return domain.NewPresentImage(img), nil
	}
}

func kvKey(name domain.AssetName) []byte {
	return []byte(fmt.Sprintf("asset:%s", name.Unwrap()))
}
