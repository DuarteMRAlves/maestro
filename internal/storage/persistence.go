package storage

import (
	"bytes"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
)

func orchestrationKey(name apitypes.OrchestrationName) []byte {
	return []byte(fmt.Sprintf("orchestration:%s", name))
}

func persistOrchestration(txn *badger.Txn, o *apitypes.Orchestration) error {
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

func loadOrchestration(o *apitypes.Orchestration, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &o.Name, &o.Phase)
	return err
}

func stageKey(name apitypes.StageName) []byte {
	return []byte(fmt.Sprintf("stage:%s", name))
}

func PersistStage(txn *badger.Txn, s *Stage) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(
		&buf,
		s.name,
		s.phase,
		s.rpcSpec.address,
		s.rpcSpec.service,
		s.rpcSpec.rpc,
		s.asset,
	)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = txn.Set(stageKey(s.Name()), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func loadStage(s *Stage, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(
		buf,
		&s.name,
		&s.phase,
		&s.rpcSpec.address,
		&s.rpcSpec.service,
		&s.rpcSpec.rpc,
		&s.asset,
	)
	return err
}

func linkKey(name apitypes.LinkName) []byte {
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
func assetKey(name apitypes.AssetName) []byte {
	return []byte(fmt.Sprintf("asset:%s", name))
}

func PersistAsset(txn *badger.Txn, a *apitypes.Asset) error {
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

func loadAsset(a *apitypes.Asset, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &a.Name, &a.Image)
	return err
}
