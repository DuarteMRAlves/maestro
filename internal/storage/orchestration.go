package storage

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/create"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
	"strings"
)

func SaveOrchestrationWithTxn(txn *badger.Txn) create.SaveOrchestration {
	return func(o domain.Orchestration) domain.OrchestrationResult {
		var (
			buf bytes.Buffer
			err error
		)
		storeStage := SaveStageWithTxn(txn)
		storeLink := SaveLinkWithTxn(txn)

		stages := o.Stages()
		links := o.Links()

		for _, s := range stages {
			res := storeStage(s)
			if res.IsError() {
				err = errdefs.PrependMsg(
					res.Error(),
					"store orchestration %s",
					o.Name(),
				)
				return domain.ErrOrchestration(err)
			}
		}
		for _, l := range links {
			res := storeLink(l)
			if res.IsError() {
				err = errdefs.PrependMsg(
					res.Error(),
					"store orchestration %s",
					o.Name(),
				)
				return domain.ErrOrchestration(err)
			}
		}
		for i, s := range stages {
			buf.WriteString(s.Name().Unwrap())
			if i+1 != len(stages) {
				// , is not a valid name char, so we can use it
				buf.WriteByte(',')
			}
		}
		buf.WriteByte(';')
		for i, l := range links {
			buf.WriteString(l.Name().Unwrap())
			if i+1 != len(links) {
				buf.WriteByte(',')
			}
		}

		err = txn.Set(orchestrationKey(o.Name()), buf.Bytes())
		if err != nil {
			err = errdefs.InternalWithMsg(
				"store orchestration %s: %v",
				o.Name(),
				err,
			)
			return domain.ErrOrchestration(err)
		}
		return domain.SomeOrchestration(o)
	}
}

func LoadOrchestrationWithTxn(txn *badger.Txn) create.LoadOrchestration {
	return func(name domain.OrchestrationName) domain.OrchestrationResult {
		var (
			item *badger.Item
			data []byte
			err  error
		)
		loadStage := LoadStageWithTxn(txn)
		loadLink := LoadLinkWithTxn(txn)

		item, err = txn.Get(orchestrationKey(name))
		if err != nil {
			err = errdefs.PrependMsg(err, "load orchestration %s", name)
			return domain.ErrOrchestration(err)
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			err = errdefs.PrependMsg(err, "load orchestration %s", name)
			return domain.ErrOrchestration(err)
		}
		buf := bytes.NewBuffer(data)
		splits := strings.Split(buf.String(), ";")
		if len(splits) != 2 {
			return domain.ErrOrchestration(
				errdefs.InternalWithMsg(
					"invalid format: expected 2 semi-colon separated values",
				),
			)
		}
		stages, err := splitToStages(loadStage, splits[0])
		if err != nil {
			err = errdefs.PrependMsg(err, "load orchestration %s", name)
			return domain.ErrOrchestration(err)
		}
		links, err := splitToLinks(loadLink, splits[1])
		if err != nil {
			err = errdefs.PrependMsg(err, "load orchestration %s", name)
			return domain.ErrOrchestration(err)
		}
		o := domain.NewOrchestration(name, stages, links)
		return domain.SomeOrchestration(o)
	}
}

func splitToStages(loadFn create.LoadStage, data string) (
	[]domain.Stage,
	error,
) {
	if len(data) == 0 {
		return []domain.Stage{}, nil
	}
	stageNames := strings.Split(data, ",")
	stages := make([]domain.Stage, 0, len(stageNames))
	for _, n := range stageNames {
		name, err := domain.NewStageName(n)
		if err != nil {
			return nil, err
		}
		res := loadFn(name)
		if res.IsError() {
			return nil, res.Error()
		}
		stages = append(stages, res.Unwrap())
	}
	return stages, nil
}

func splitToLinks(loadFn create.LoadLink, data string) (
	[]domain.Link,
	error,
) {
	if len(data) == 0 {
		return []domain.Link{}, nil
	}
	linksNames := strings.Split(data, ",")
	links := make([]domain.Link, 0, len(linksNames))
	for _, n := range linksNames {
		name, err := domain.NewLinkName(n)
		if err != nil {
			return nil, err
		}
		res := loadFn(name)
		if res.IsError() {
			return nil, res.Error()
		}
		links = append(links, res.Unwrap())
	}
	return links, nil
}

func orchestrationKey(name domain.OrchestrationName) []byte {
	return []byte(fmt.Sprintf("orchestration:%s", name.Unwrap()))
}
