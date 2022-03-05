package link

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/types"
	"github.com/dgraph-io/badger/v3"
	"strings"
)

func StoreWithTxn(txn *badger.Txn) func(types.Link) types.LinkResult {
	return func(l types.Link) types.LinkResult {
		var (
			buf bytes.Buffer
			err error
		)
		storeStage := stage.StoreWithTxn(txn)

		sourceStage := l.Source().Stage()
		targetStage := l.Target().Stage()

		sourceRes := storeStage(sourceStage)
		if sourceRes.IsError() {
			err = sourceRes.Error()
			err = errdefs.PrependMsg(err, "store link: %s", l.Name())
			return Err(err)
		}
		targetRes := storeStage(targetStage)
		if targetRes.IsError() {
			err = sourceRes.Error()
			err = errdefs.PrependMsg(err, "store link %s", l.Name())
			return Err(err)
		}
		linkPersistenceInfoToBuf(&buf, l)
		err = txn.Set(kvKey(l.Name()), buf.Bytes())
		if err != nil {
			err = errdefs.InternalWithMsg("store link %s: %v", l.Name(), err)
			return Err(err)
		}
		return Some(l)
	}
}

func LoadWithTxn(txn *badger.Txn) func(types.LinkName) types.LinkResult {
	return func(name types.LinkName) types.LinkResult {
		var (
			item *badger.Item
			data []byte
			err  error
		)
		item, err = txn.Get(kvKey(name))
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load link %s", name))
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load link %s", name))
		}
		buf := bytes.NewBuffer(data)
		splits := strings.Split(buf.String(), ";")
		if len(splits) != 4 {
			return Err(
				errdefs.InternalWithMsg(
					"invalid format: expected 4 semi-colon separated values",
				),
			)
		}
		loadStage := stage.LoadWithTxn(txn)
		source, err := loadEndpoint(loadStage, splits[0], splits[1])
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load link %s", name))
		}
		target, err := loadEndpoint(loadStage, splits[2], splits[3])
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load link %s", name))
		}
		return Some(NewLink(name, source, target))
	}
}

func linkPersistenceInfoToBuf(buf *bytes.Buffer, l types.Link) {
	endpointToBuf(buf, l.Source())
	buf.WriteByte(';')
	endpointToBuf(buf, l.Target())
}

func endpointToBuf(buf *bytes.Buffer, e types.LinkEndpoint) {
	buf.WriteString(e.Stage().Name().Unwrap())
	buf.WriteByte(';')
	if e.Field().Present() {
		buf.WriteString(e.Field().Unwrap().Unwrap())
	}
}

func loadEndpoint(
	loadStage func(types.StageName) types.StageResult,
	nameData string,
	fieldData string,
) (types.LinkEndpoint, error) {
	name, err := stage.NewStageName(nameData)
	if err != nil {
		return nil, err
	}
	res := loadStage(name)
	if res.IsError() {
		return nil, res.Error()
	}
	stage := res.Unwrap()
	field, err := loadField(fieldData)
	if err != nil {
		return nil, err
	}
	return NewLinkEndpoint(stage, field), nil
}

func loadField(data string) (types.OptionalMessageField, error) {
	if data == "" {
		return NewEmptyMessageField(), nil
	} else {
		f, err := NewMessageField(data)
		if err != nil {
			return nil, err
		}
		return NewPresentMessageField(f), nil
	}
}

func kvKey(name types.LinkName) []byte {
	return []byte(fmt.Sprintf("link:%s", name.Unwrap()))
}
