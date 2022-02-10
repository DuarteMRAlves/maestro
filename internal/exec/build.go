package exec

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc"
	"time"
)

// Builder constructs an execution from an orchestration.
type Builder struct {
	orchestration *api.Orchestration
	stageMap      *StageMap

	orchestrationName api.OrchestrationName
	stages            map[api.StageName]*api.Stage
	links             map[api.LinkName]*api.Link
	rpcs              map[api.StageName]rpc.RPC
	inputs            map[api.StageName]*InputDesc
	outputs           map[api.StageName]*OutputDesc

	txnHelper  *kv.TxnHelper
	rpcManager rpc.Manager
}

func newBuilder(txn *badger.Txn, rpcManager rpc.Manager) *Builder {
	return &Builder{
		stageMap:   NewStageMap(),
		txnHelper:  kv.NewTxnHelper(txn),
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
	err = b.buildExecStages()
	if err != nil {
		return nil, err
	}
	e := &Execution{
		orchestration: b.orchestration,
		stages:        b.stageMap,
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
	b.inputs = make(map[api.StageName]*InputDesc, len(b.stages))
	b.outputs = make(map[api.StageName]*OutputDesc, len(b.stages))
	for name := range b.stages {
		stageRpc, ok := b.rpcs[name]
		if !ok {
			return errdefs.InternalWithMsg("rpc not found for %s", name)
		}
		b.inputs[name] = NewInputDesc().WithMessage(stageRpc.Input())
		b.outputs[name] = NewOutputBuilder()
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

		conn := NewLink(l)

		err := b.inputs[targetName].WithConnection(conn)
		if err != nil {
			return errdefs.PrependMsg(
				err,
				"register %s in %s input",
				l.Name,
				targetName,
			)
		}

		err = b.outputs[sourceName].WithConnection(conn)
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

func (b *Builder) buildExecStages() error {
	var (
		inputChan, outputChan chan *State
		stage                 Stage
		err                   error
	)

	for name, apiStage := range b.stages {

		stageRpc, ok := b.rpcs[name]
		if !ok {
			return errdefs.InternalWithMsg("rpc not found for %s", name)
		}
		inputBuilder, ok := b.inputs[name]
		if !ok {
			return errdefs.InternalWithMsg(
				"input builder not found for %s",
				name,
			)
		}
		outputBuilder, ok := b.outputs[name]
		if !ok {
			return errdefs.InternalWithMsg(
				"output builder not found for %s",
				name,
			)
		}
		inputChan, stage, err = inputBuilder.BuildExecutionResources()
		if err != nil {
			return errdefs.PrependMsg(err, "input build error for %s", name)
		}
		if stage != nil {
			b.stageMap.AddInputStage(name, stage)
		}
		outputChan, stage, err = outputBuilder.BuildExecutionResources()
		if err != nil {
			return errdefs.PrependMsg(err, "output build error for %s", name)
		}
		if stage != nil {
			b.stageMap.AddOutputStage(name, stage)
		}
		cfg := &StageCfg{
			Address: apiStage.Address,
			Rpc:     stageRpc,
			Input:   inputChan,
			Output:  outputChan,
		}

		stage, err = NewStage(cfg)
		if err != nil {
			return errdefs.PrependMsg(err, "build rpc stage")
		}
		b.stageMap.AddRpcStage(name, stage)
	}
	return nil
}

// InputDesc stores the necessary information to handle the input of a stage.
type InputDesc struct {
	connections []*Link
	msg         rpc.Message
}

func NewInputDesc() *InputDesc {
	return &InputDesc{
		connections: []*Link{},
	}
}

func (i *InputDesc) WithMessage(msg rpc.Message) *InputDesc {
	i.msg = msg
	return i
}

func (i *InputDesc) WithConnection(c *Link) error {
	// A previous link that consumes the entire message already exists
	if len(i.connections) == 1 && i.connections[0].HasEmptyTargetField() {
		return errdefs.FailedPreconditionWithMsg(
			"link that receives the full message already exists",
		)
	}
	for _, prev := range i.connections {
		if prev.HasSameTargetField(c) {
			return errdefs.InvalidArgumentWithMsg(
				"link with the same target field already registered: %s",
				prev.LinkName(),
			)
		}
	}
	i.connections = append(i.connections, c)
	return nil
}

func (i *InputDesc) BuildExecutionResources() (chan *State, Stage, error) {
	switch len(i.connections) {
	case 0:
		if i.msg == nil {
			return nil, nil, errdefs.FailedPreconditionWithMsg(
				"message required without 0 connections",
			)
		}
		ch := make(chan *State)
		s := NewSourceStage(1, ch, i.msg)
		return ch, s, nil
	case 1:
		return i.connections[0].Chan(), nil, nil
	default:
		return nil, nil, errdefs.FailedPreconditionWithMsg(
			"too many connections: expected 0 or 1 but received %d",
			len(i.connections),
		)
	}
}

// OutputDesc registers the several connections for an output.
type OutputDesc struct {
	connections []*Link
}

func NewOutputBuilder() *OutputDesc {
	return &OutputDesc{
		connections: []*Link{},
	}
}

func (o *OutputDesc) WithConnection(c *Link) error {
	for _, prev := range o.connections {
		if prev.HasSameLinkName(c) {
			return errdefs.InvalidArgumentWithMsg(
				"Link with an equal name already registered: %s",
				prev.LinkName(),
			)
		}
	}

	o.connections = append(o.connections, c)
	return nil
}

func (o *OutputDesc) BuildExecutionResources() (chan *State, Stage, error) {
	switch len(o.connections) {
	case 0:
		ch := make(chan *State)
		s := NewSinkOutput(ch)
		return ch, s, nil
	case 1:
		return o.connections[0].Chan(), nil, nil
	default:
		return nil, nil, errdefs.FailedPreconditionWithMsg(
			"too many connections: expected 0 or 1 but received %d",
			len(o.connections),
		)
	}
}
