package storage

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
)

const (
	orchestrationPrefix = "orchestration:"
	stagePrefix         = "stage:"
	linkPrefix          = "link:"
	assetPrefix         = "asset:"
)

// TxnHelper wraps a transaction and offers utility functions for common
// transactional operations.
type TxnHelper struct {
	txn *badger.Txn
}

func NewTxnHelper(txn *badger.Txn) *TxnHelper {
	return &TxnHelper{txn: txn}
}

type OrchestrationVisitor = func(*api.Orchestration) error
type StageVisitor = func(*api.Stage) error
type LinkVisitor = func(*api.Link) error
type AssetVisitor = func(*api.Asset) error

type IterVisitor = func([]byte) error

type IterOpts struct {
	copy       bool
	badgerOpts badger.IteratorOptions
}

func DefaultIterOpts() IterOpts {
	return IterOpts{
		copy:       true,
		badgerOpts: badger.DefaultIteratorOptions,
	}
}

func (h *TxnHelper) SaveOrchestration(o *api.Orchestration) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(&buf, o.Name, o.Phase)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = h.txn.Set(orchestrationKey(o.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) LoadOrchestration(
	o *api.Orchestration,
	name api.OrchestrationName,
) error {
	var (
		item *badger.Item
		data []byte
		err  error
	)
	item, err = h.txn.Get(orchestrationKey(name))
	if err != nil {
		return err
	}
	data, err = item.ValueCopy(nil)
	return loadOrchestration(o, data)
}

func (h *TxnHelper) IterOrchestrations(
	vis OrchestrationVisitor,
	opts IterOpts,
) error {
	var o api.Orchestration
	iterVis := func(data []byte) error {
		err := loadOrchestration(&o, data)
		if err != nil {
			return err
		}
		return vis(&o)
	}
	return h.iterValues(iterVis, opts, []byte(orchestrationPrefix))
}

func loadOrchestration(o *api.Orchestration, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &o.Name, &o.Phase)
	return err
}

func (h *TxnHelper) SaveStage(s *api.Stage) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(
		&buf,
		s.Name,
		s.Phase,
		s.Service,
		s.Rpc,
		s.Address,
		s.Asset,
	)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = h.txn.Set(stageKey(s.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) LoadStage(s *api.Stage, name api.StageName) error {
	var (
		item *badger.Item
		data []byte
		err  error
	)
	item, err = h.txn.Get(stageKey(name))
	if err != nil {
		return err
	}
	data, err = item.ValueCopy(nil)
	return loadStage(s, data)
}

func (h *TxnHelper) IterStages(
	vis StageVisitor,
	opts IterOpts,
) error {
	var s api.Stage
	iterVis := func(data []byte) error {
		err := loadStage(&s, data)
		if err != nil {
			return err
		}
		return vis(&s)
	}
	return h.iterValues(iterVis, opts, []byte(stagePrefix))
}

func loadStage(s *api.Stage, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(
		buf,
		&s.Name,
		&s.Phase,
		&s.Service,
		&s.Rpc,
		&s.Address,
		&s.Asset,
	)
	return err
}

func (h *TxnHelper) SaveLink(l *api.Link) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(
		&buf,
		l.Name,
		l.SourceStage,
		l.SourceField,
		l.TargetStage,
		l.TargetField,
	)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = h.txn.Set(linkKey(l.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) IterLinks(
	vis LinkVisitor,
	opts IterOpts,
) error {
	var l api.Link
	iterVis := func(data []byte) error {
		err := loadLink(&l, data)
		if err != nil {
			return err
		}
		return vis(&l)
	}
	return h.iterValues(iterVis, opts, []byte(linkPrefix))
}

func loadLink(l *api.Link, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(
		buf,
		&l.Name,
		&l.SourceStage,
		&l.SourceField,
		&l.TargetStage,
		&l.TargetField,
	)
	return err
}

func (h *TxnHelper) SaveAsset(a *api.Asset) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(&buf, a.Name, a.Image)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = h.txn.Set(assetKey(a.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) IterAssets(
	vis AssetVisitor,
	opts IterOpts,
) error {
	var a api.Asset
	iterVis := func(data []byte) error {
		err := loadAsset(&a, data)
		if err != nil {
			return err
		}
		return vis(&a)
	}
	return h.iterValues(iterVis, opts, []byte(assetPrefix))
}

func loadAsset(a *api.Asset, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &a.Name, &a.Image)
	return err
}

func (h *TxnHelper) iterValues(
	vis IterVisitor,
	opts IterOpts,
	prefix []byte,
) error {
	var (
		visitItemFn func(IterVisitor, *badger.Item) error
		err         error
	)
	it := h.txn.NewIterator(opts.badgerOpts)
	defer it.Close()

	if opts.copy {
		visitItemFn = visitWithCopy
	} else {
		visitItemFn = visitWithoutCopy
	}

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		err = visitItemFn(vis, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func visitWithCopy(vis IterVisitor, item *badger.Item) error {
	var (
		cp  []byte
		err error
	)
	if cp, err = item.ValueCopy(nil); err != nil {
		return errdefs.InternalWithMsg("value copy: %v", err)
	}
	return vis(cp)
}

func visitWithoutCopy(vis IterVisitor, item *badger.Item) error {
	return item.Value(vis)
}

func orchestrationKey(name api.OrchestrationName) []byte {
	return []byte(fmt.Sprintf("%s%s", orchestrationPrefix, name))
}

func stageKey(name api.StageName) []byte {
	return []byte(fmt.Sprintf("%s%s", stagePrefix, name))
}

func linkKey(name api.LinkName) []byte {
	return []byte(fmt.Sprintf("%s%s", linkPrefix, name))
}

func assetKey(name api.AssetName) []byte {
	return []byte(fmt.Sprintf("%s%s", assetPrefix, name))
}
