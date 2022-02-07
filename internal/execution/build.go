package execution

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc"
	"time"
)

// Builder constructs an execution from an orchestration.
type Builder struct {
	orchestration *api.Orchestration
	workers       map[api.StageName]Worker

	orchestrationName api.OrchestrationName
	stages            map[api.StageName]*api.Stage
	links             map[api.LinkName]*api.Link
	rpcs              map[api.StageName]rpc.RPC
	inputBuilders     map[api.StageName]*InputBuilder
	outputBuilders    map[api.StageName]*OutputBuilder

	txnHelper  *storage.TxnHelper
	rpcManager rpc.Manager
}

func newBuilder(txn *badger.Txn, rpcManager rpc.Manager) *Builder {
	return &Builder{
		txnHelper:  storage.NewTxnHelper(txn),
		rpcManager: rpcManager,
	}
}

func (b *Builder) withOrchestration(name api.OrchestrationName) *Builder {
	b.orchestrationName = name
	return b
}

func (b *Builder) build() (*Execution, error) {
	var err error
	b.orchestration = &api.Orchestration{}
	err = b.txnHelper.LoadOrchestration(b.orchestration, b.orchestrationName)
	if err != nil {
		return nil, err
	}
	err = b.loadStages()
	if err != nil {
		return nil, err
	}
	err = b.loadLinks()
	if err != nil {
		return nil, err
	}
	err = b.loadRpcs()
	if err != nil {
		return nil, err
	}
	err = b.loadInputsAndOutputs()
	if err != nil {
		return nil, err
	}
	err = b.buildWorkers()
	if err != nil {
		return nil, err
	}
	e := &Execution{
		orchestration: b.orchestration,
		workers:       b.workers,
	}
	return e, nil
}

func (b *Builder) loadStages() error {
	var (
		loaded *api.Stage
		err    error
	)

	o := b.orchestration

	b.stages = make(map[api.StageName]*api.Stage, len(o.Stages))
	for _, s := range o.Stages {
		loaded = &api.Stage{}
		err = b.txnHelper.LoadStage(loaded, s)
		if err != nil {
			return errdefs.PrependMsg(err, "load stages")
		}
		b.stages[s] = loaded
	}
	return nil
}

func (b *Builder) loadLinks() error {
	var (
		loaded *api.Link
		err    error
	)

	o := b.orchestration

	b.links = make(map[api.LinkName]*api.Link, len(o.Links))
	for _, l := range o.Links {
		loaded = &api.Link{}
		err = b.txnHelper.LoadLink(loaded, l)
		if err != nil {
			return errdefs.PrependMsg(err, "load links")
		}
		b.links[l] = loaded
	}
	return nil
}

func (b *Builder) loadRpcs() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b.rpcs = make(map[api.StageName]rpc.RPC, len(b.stages))
	for name, s := range b.stages {
		stageRpc, err := b.loadRpc(ctx, s)
		if err != nil {
			return errdefs.PrependMsg(err, "load rpcs")
		}
		b.rpcs[name] = stageRpc
	}
	return nil
}

func (b *Builder) loadRpc(ctx context.Context, s *api.Stage) (rpc.RPC, error) {
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		return nil, errdefs.InternalWithMsg(
			"dial %s: %v",
			s.Name,
			err,
		)
	}
	defer conn.Close()
	return b.rpcManager.GetRpc(ctx, conn, s)
}

func (b *Builder) loadInputsAndOutputs() error {
	b.inputBuilders = make(map[api.StageName]*InputBuilder, len(b.stages))
	b.outputBuilders = make(map[api.StageName]*OutputBuilder, len(b.stages))
	for name := range b.stages {
		stageRpc, ok := b.rpcs[name]
		if !ok {
			return errdefs.InternalWithMsg("rpc not found for %s", name)
		}
		b.inputBuilders[name] = NewInputBuilder().WithMessage(stageRpc.Input())
		b.outputBuilders[name] = NewOutputBuilder()
	}

	for _, l := range b.links {
		sourceName, targetName := l.SourceStage, l.TargetStage

		sourceRpc, ok := b.rpcs[sourceName]
		if !ok {
			return errdefs.InternalWithMsg("rpc not found for %s", sourceName)
		}
		targetRpc, ok := b.rpcs[targetName]
		if !ok {
			return errdefs.InternalWithMsg("rpc not found for %s", targetName)
		}

		sourceMsg := sourceRpc.Output()
		if l.SourceField != "" {
			sourceMsg, ok = sourceMsg.GetMessageField(l.SourceField)
			if !ok {
				return errdefs.NotFoundWithMsg(
					"field %s not found for source message %s",
					l.SourceField,
					sourceMsg.FullyQualifiedName(),
				)
			}
		}
		targetMsg := targetRpc.Input()
		if l.TargetField != "" {
			targetMsg, ok = targetMsg.GetMessageField(l.TargetField)
			if !ok {
				return errdefs.NotFoundWithMsg(
					"field %s not found for target message %s",
					l.TargetField,
					targetMsg.FullyQualifiedName(),
				)
			}
		}

		if !sourceMsg.Compatible(targetMsg) {
			return errdefs.InvalidArgumentWithMsg(
				"incompatible messages for link %s: source is %s, target is %s",
				l.Name,
				sourceMsg.FullyQualifiedName(),
				targetMsg.FullyQualifiedName(),
			)
		}

		conn, err := NewConnection(l)
		if err != nil {
			return errdefs.PrependMsg(
				err,
				"create internal connection for %s",
				l.Name,
			)
		}

		err = b.inputBuilders[targetName].WithConnection(conn)
		if err != nil {
			return errdefs.PrependMsg(
				err,
				"register %s in %s input",
				l.Name,
				targetName,
			)
		}

		err = b.outputBuilders[sourceName].WithConnection(conn)
		if err != nil {
			return errdefs.PrependMsg(
				err,
				"register %s in %s output",
				l.Name,
				targetName,
			)
		}
	}
	return nil
}

func (b *Builder) buildWorkers() error {
	var (
		input  Input
		output Output
		err    error
	)

	b.workers = make(map[api.StageName]Worker, len(b.stages))
	sources := make([]api.StageName, 0)
	sinks := make([]api.StageName, 0)
	for name, stage := range b.stages {

		stageRpc, ok := b.rpcs[name]
		if !ok {
			return errdefs.InternalWithMsg("rpc not found for %s", name)
		}
		inputBuilder, ok := b.inputBuilders[name]
		if !ok {
			return errdefs.InternalWithMsg(
				"input builder not found for %s",
				name,
			)
		}
		outputBuilder, ok := b.outputBuilders[name]
		if !ok {
			return errdefs.InternalWithMsg(
				"output builder not found for %s",
				name,
			)
		}
		input, err = inputBuilder.Build()
		if err != nil {
			return errdefs.PrependMsg(err, "input build error for %s", name)
		}
		output, err = outputBuilder.Build()
		if err != nil {
			return errdefs.PrependMsg(err, "output build error for %s", name)
		}
		cfg := &WorkerCfg{
			Address: stage.Address,
			Rpc:     stageRpc,
			Input:   input,
			Output:  output,
		}

		if cfg.Input.IsSource() {
			sources = append(sources, name)
		}
		if cfg.Output.IsSink() {
			sinks = append(sinks, name)
		}

		b.workers[name], err = NewWorker(cfg)
		if err != nil {
			return errdefs.PrependMsg(err, "build workers")
		}
	}
	if len(sources) != 1 {
		return errdefs.InvalidArgumentWithMsg(
			"expected one source stage but found %d: %v",
			len(sources),
			sources,
		)
	}
	if len(sinks) != 1 {
		return errdefs.InvalidArgumentWithMsg(
			"expected one sink stage but found %d: %v",
			len(sinks),
			sinks,
		)
	}
	return nil
}
