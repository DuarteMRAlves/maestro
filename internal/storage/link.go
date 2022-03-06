package storage

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/create"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	"github.com/dgraph-io/badger/v3"
	"strings"
)

func SaveLinkWithTxn(txn *badger.Txn) create.SaveLink {
	return func(l execute.Link) execute.LinkResult {
		var (
			buf bytes.Buffer
			err error
		)
		storeStage := SaveStageWithTxn(txn)

		sourceStage := l.Source().Stage()
		targetStage := l.Target().Stage()

		sourceRes := storeStage(sourceStage)
		if sourceRes.IsError() {
			err = sourceRes.Error()
			err = errdefs.PrependMsg(err, "store link: %s", l.Name())
			return execute.ErrLink(err)
		}
		targetRes := storeStage(targetStage)
		if targetRes.IsError() {
			err = sourceRes.Error()
			err = errdefs.PrependMsg(err, "store link %s", l.Name())
			return execute.ErrLink(err)
		}
		linkPersistenceInfoToBuf(&buf, l)
		err = txn.Set(linkKey(l.Name()), buf.Bytes())
		if err != nil {
			err = errdefs.InternalWithMsg("store link %s: %v", l.Name(), err)
			return execute.ErrLink(err)
		}
		return execute.SomeLink(l)
	}
}

func LoadLinkWithTxn(txn *badger.Txn) execute.LoadLink {
	return func(name domain.LinkName) execute.LinkResult {
		var (
			item *badger.Item
			data []byte
			err  error
		)
		item, err = txn.Get(linkKey(name))
		if err != nil {
			err = errdefs.PrependMsg(err, "load link %s", name)
			return execute.ErrLink(err)
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			err = errdefs.PrependMsg(err, "load link %s", name)
			return execute.ErrLink(err)
		}
		buf := bytes.NewBuffer(data)
		splits := strings.Split(buf.String(), ";")
		if len(splits) != 4 {
			return execute.ErrLink(
				errdefs.InternalWithMsg(
					"invalid format: expected 4 semi-colon separated values",
				),
			)
		}
		loadStage := LoadStageWithTxn(txn)
		source, err := loadEndpoint(loadStage, splits[0], splits[1])
		if err != nil {
			err = errdefs.PrependMsg(err, "load link %s", name)
			return execute.ErrLink(err)
		}
		target, err := loadEndpoint(loadStage, splits[2], splits[3])
		if err != nil {
			err = errdefs.PrependMsg(err, "load link %s", name)
			return execute.ErrLink(err)
		}
		return execute.SomeLink(execute.NewLink(name, source, target))
	}
}

func linkPersistenceInfoToBuf(buf *bytes.Buffer, l execute.Link) {
	endpointToBuf(buf, l.Source())
	buf.WriteByte(';')
	endpointToBuf(buf, l.Target())
}

func endpointToBuf(buf *bytes.Buffer, e execute.LinkEndpoint) {
	buf.WriteString(e.Stage().Name().Unwrap())
	buf.WriteByte(';')
	if e.Field().Present() {
		buf.WriteString(e.Field().Unwrap().Unwrap())
	}
}

func loadEndpoint(
	loadStage func(domain.StageName) domain.StageResult,
	nameData string,
	fieldData string,
) (execute.LinkEndpoint, error) {
	name, err := domain.NewStageName(nameData)
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
	return execute.NewLinkEndpoint(stage, field), nil
}

func loadField(data string) (domain.OptionalMessageField, error) {
	if data == "" {
		return domain.NewEmptyMessageField(), nil
	} else {
		f, err := domain.NewMessageField(data)
		if err != nil {
			return nil, err
		}
		return domain.NewPresentMessageField(f), nil
	}
}

func linkKey(name domain.LinkName) []byte {
	return []byte(fmt.Sprintf("link:%s", name.Unwrap()))
}
