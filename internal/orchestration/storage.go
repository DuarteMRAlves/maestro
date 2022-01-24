package orchestration

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

func persistOrchestration(txn *badger.Txn, o *Orchestration) error {
	var (
		buf bytes.Buffer
		err error
	)
	_, err = fmt.Fprintln(&buf, o.name, o.phase)
	if err != nil {
		return errdefs.InternalWithMsg("encoding error: %v", err)
	}
	err = txn.Set(orchestrationKey(o.name), buf.Bytes())
	if err != nil {
		return errdefs.InternalWithMsg("storage error: %v", err)
	}
	return nil
}

func loadOrchestration(o *Orchestration, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(buf, &o.name, &o.phase)
	return err
}
