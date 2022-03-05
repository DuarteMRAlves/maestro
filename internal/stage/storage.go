package stage

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
	"strings"
)

func StoreWithTxn(txn *badger.Txn) func(domain.Stage) domain.StageResult {
	return func(s domain.Stage) domain.StageResult {
		var (
			buf bytes.Buffer
			err error
		)
		stageToBuf(&buf, s)
		err = txn.Set(kvKey(s.Name()), buf.Bytes())
		if err != nil {
			return Err(errdefs.InternalWithMsg("storage error: %v", err))
		}
		return Some(s)
	}
}

func LoadWithTxn(txn *badger.Txn) func(domain.StageName) domain.StageResult {
	return func(name domain.StageName) domain.StageResult {
		var (
			item *badger.Item
			data []byte
			err  error
		)
		item, err = txn.Get(kvKey(name))
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load stage %s", name))
		}
		data, err = item.ValueCopy(nil)
		buf := bytes.NewBuffer(data)
		splits := strings.Split(buf.String(), ";")
		if len(splits) != 3 {
			return Err(
				errdefs.InternalWithMsg(
					"invalid format: expected 3 semi-colon separated values",
				),
			)
		}
		a, err := stringToAddress(splits[0])
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load stage %s", name))
		}
		s, err := stringToService(splits[1])
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load stage %s", name))
		}
		m, err := stringToMethod(splits[2])
		if err != nil {
			return Err(errdefs.PrependMsg(err, "load stage %s", name))
		}
		methodCtx := NewMethodContext(a, s, m)
		return Some(NewStage(name, methodCtx))
	}
}

func stageToBuf(buf *bytes.Buffer, s domain.Stage) {
	methodCtxToBuf(buf, s.MethodContext())
}

func methodCtxToBuf(buf *bytes.Buffer, m domain.MethodContext) {
	buf.WriteString(m.Address().Unwrap())
	buf.WriteByte(';')
	if m.Service().Present() {
		buf.WriteString(m.Service().Unwrap().Unwrap())
	}
	buf.WriteByte(';')
	if m.Method().Present() {
		buf.WriteString(m.Method().Unwrap().Unwrap())
	}
}

func stringToAddress(data string) (domain.Address, error) {
	return NewAddress(data)
}

func stringToService(data string) (domain.OptionalService, error) {
	if data == "" {
		return NewEmptyService(), nil
	} else {
		s, err := NewService(data)
		if err != nil {
			return nil, err
		}
		return NewPresentService(s), nil
	}
}

func stringToMethod(data string) (domain.OptionalMethod, error) {
	if data == "" {
		return NewEmptyMethod(), nil
	} else {
		m, err := NewMethod(data)
		if err != nil {
			return nil, err
		}
		return NewPresentMethod(m), nil
	}
}

func kvKey(name domain.StageName) []byte {
	return []byte(fmt.Sprintf("stage:%s", name.Unwrap()))
}
