package storage

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
)

func orchestrationKey(name api.OrchestrationName) []byte {
	return []byte(fmt.Sprintf("orchestration:%s", name))
}

func persistOrchestration(txn *badger.Txn, o *api.Orchestration) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(&buf, o.Name, o.Phase)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = txn.Set(orchestrationKey(o.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func loadOrchestration(o *api.Orchestration, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &o.Name, &o.Phase)
	return err
}

func stageKey(name api.StageName) []byte {
	return []byte(fmt.Sprintf("stage:%s", name))
}

func PersistStage(txn *badger.Txn, s *api.Stage) error {
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
	err = txn.Set(stageKey(s.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
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

func linkKey(name api.LinkName) []byte {
	return []byte(fmt.Sprintf("link:%s", name))
}

func PersistLink(txn *badger.Txn, l *Link) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(
		&buf,
		l.name,
		l.sourceStage,
		l.sourceField,
		l.targetStage,
		l.targetField,
	)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = txn.Set(linkKey(l.name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func loadLink(l *Link, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(
		buf,
		&l.name,
		&l.sourceStage,
		&l.sourceField,
		&l.targetStage,
		&l.targetField,
	)
	return err
}

// assetKey returns the image key for an asset with the given name
func assetKey(name api.AssetName) []byte {
	return []byte(fmt.Sprintf("asset:%s", name))
}

func PersistAsset(txn *badger.Txn, a *api.Asset) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(&buf, a.Name, a.Image)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = txn.Set(assetKey(a.Name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func loadAsset(a *api.Asset, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &a.Name, &a.Image)
	return err
}
