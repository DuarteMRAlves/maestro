package storage

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
	"strings"
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
	buf.WriteString(string(o.Name))
	// ; is not a valid name char, so we can use it
	buf.WriteByte(';')
	buf.WriteString(string(o.Phase))
	buf.WriteByte(';')
	for i, s := range o.Stages {
		buf.WriteString(string(s))
		if i+1 != len(o.Stages) {
			// , is not a valid name char, so we can use it
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(';')
	for i, l := range o.Links {
		buf.WriteString(string(l))
		if i+1 != len(o.Links) {
			buf.WriteByte(',')
		}
	}
	err = h.txn.Set(orchestrationKey(o.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) ContainsOrchestration(name api.OrchestrationName) bool {
	item, _ := h.txn.Get(orchestrationKey(name))
	return item != nil
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
	splits := strings.Split(buf.String(), ";")
	if len(splits) != 4 {
		return errdefs.InternalWithMsg(
			"invalid format: expected 4 semi-colon separated values",
		)
	}
	o.Name = api.OrchestrationName(splits[0])
	o.Phase = api.OrchestrationPhase(splits[1])

	// Stages are empty
	if len(splits[2]) == 0 {
		o.Stages = []api.StageName{}
	} else {
		stages := strings.Split(splits[2], ",")
		o.Stages = make([]api.StageName, 0, len(stages))
		for _, s := range stages {
			o.Stages = append(o.Stages, api.StageName(s))
		}
	}

	// Links are empty
	if len(splits[3]) == 0 {
		o.Links = []api.LinkName{}
	} else {
		links := strings.Split(splits[3], ",")
		o.Links = make([]api.LinkName, 0, len(links))
		for _, l := range links {
			o.Links = append(o.Links, api.LinkName(l))
		}
	}
	return nil
}

func (h *TxnHelper) SaveStage(s *api.Stage) error {
	var (
		buf bytes.Buffer
		err error
	)
	buf.WriteString(string(s.Name))
	buf.WriteByte(';')
	buf.WriteString(string(s.Phase))
	buf.WriteByte(';')
	buf.WriteString(s.Service)
	buf.WriteByte(';')
	buf.WriteString(s.Rpc)
	buf.WriteByte(';')
	buf.WriteString(s.Address)
	buf.WriteByte(';')
	buf.WriteString(string(s.Orchestration))
	buf.WriteByte(';')
	buf.WriteString(string(s.Asset))
	err = h.txn.Set(stageKey(s.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) ContainsStage(name api.StageName) bool {
	item, _ := h.txn.Get(stageKey(name))
	return item != nil
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
	splits := strings.Split(buf.String(), ";")
	if len(splits) != 7 {
		return errdefs.InternalWithMsg(
			"invalid format: expected 7 semi-colon separated values",
		)
	}
	s.Name = api.StageName(splits[0])
	s.Phase = api.StagePhase(splits[1])
	s.Service = splits[2]
	s.Rpc = splits[3]
	s.Address = splits[4]
	s.Orchestration = api.OrchestrationName(splits[5])
	s.Asset = api.AssetName(splits[6])
	return nil
}

func (h *TxnHelper) SaveLink(l *api.Link) error {
	var (
		buf bytes.Buffer
		err error
	)
	buf.WriteString(string(l.Name))
	buf.WriteByte(';')
	buf.WriteString(string(l.SourceStage))
	buf.WriteByte(';')
	buf.WriteString(l.SourceField)
	buf.WriteByte(';')
	buf.WriteString(string(l.TargetStage))
	buf.WriteByte(';')
	buf.WriteString(l.TargetField)
	buf.WriteByte(';')
	buf.WriteString(string(l.Orchestration))

	err = h.txn.Set(linkKey(l.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) ContainsLink(name api.LinkName) bool {
	item, _ := h.txn.Get(linkKey(name))
	return item != nil
}

func (h *TxnHelper) LoadLink(l *api.Link, name api.LinkName) error {
	var (
		item *badger.Item
		data []byte
		err  error
	)
	item, err = h.txn.Get(linkKey(name))
	if err != nil {
		return err
	}
	data, err = item.ValueCopy(nil)
	return loadLink(l, data)
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
	splits := strings.Split(buf.String(), ";")
	if len(splits) != 6 {
		return errdefs.InternalWithMsg(
			"invalid format: expected 6 semi-colon separated values",
		)
	}
	l.Name = api.LinkName(splits[0])
	l.SourceStage = api.StageName(splits[1])
	l.SourceField = splits[2]
	l.TargetStage = api.StageName(splits[3])
	l.TargetField = splits[4]
	l.Orchestration = api.OrchestrationName(splits[5])
	return nil
}

func (h *TxnHelper) SaveAsset(a *api.Asset) error {
	var (
		buf bytes.Buffer
		err error
	)
	buf.WriteString(string(a.Name))
	buf.WriteByte(';')
	buf.WriteString(a.Image)
	err = h.txn.Set(assetKey(a.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func (h *TxnHelper) ContainsAsset(name api.AssetName) bool {
	item, _ := h.txn.Get(assetKey(name))
	return item != nil
}

func (h *TxnHelper) LoadAsset(a *api.Asset, name api.AssetName) error {
	var (
		item *badger.Item
		data []byte
		err  error
	)
	item, err = h.txn.Get(assetKey(name))
	if err != nil {
		return err
	}
	data, err = item.ValueCopy(nil)
	return loadAsset(a, data)
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
	splits := strings.Split(buf.String(), ";")
	if len(splits) != 2 {
		return errdefs.InternalWithMsg(
			"invalid format: expected 2 semi-colon separated values",
		)
	}
	a.Name = api.AssetName(splits[0])
	a.Image = splits[1]
	return nil
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
