package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
)

// CreateStageContext offers a context for the creation of a stage.
type CreateStageContext struct {
	// req contains the received request. This request should not be changed.
	req *api.CreateStageRequest

	txnHelper *TxnHelper

	orchestration   *api.Orchestration
	inferredAddress string
}

func newCreateStageContext(
	req *api.CreateStageRequest,
	txnHelper *TxnHelper,
) *CreateStageContext {
	return &CreateStageContext{
		req:       req,
		txnHelper: txnHelper,
	}
}

func (c *CreateStageContext) validateAndComplete() error {
	var orchestrationName api.OrchestrationName
	if ok, err := util.ArgNotNil(c.req, "req"); !ok {
		return err
	}
	if !IsValidStageName(c.req.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			c.req.Name,
		)
	}
	if c.txnHelper.ContainsStage(c.req.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"stage '%v' already exists",
			c.req.Name,
		)
	}
	orchestrationName = c.req.Orchestration
	if orchestrationName == "" {
		orchestrationName = DefaultOrchestrationName
	}
	if !c.txnHelper.ContainsOrchestration(orchestrationName) {
		return errdefs.NotFoundWithMsg(
			"orchestration '%v' not found",
			orchestrationName,
		)
	}
	c.orchestration = &api.Orchestration{}
	err := c.txnHelper.LoadOrchestration(c.orchestration, orchestrationName)
	if err != nil {
		return errdefs.PrependMsg(err, "Unable to load orchestration")
	}
	// Asset is not required but if specified should exist.
	if c.req.Asset != "" && !c.txnHelper.ContainsAsset(c.req.Asset) {
		return errdefs.NotFoundWithMsg(
			"asset '%v' not found",
			c.req.Asset,
		)
	}
	if c.req.Address != "" && c.req.Host != "" {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and host for stage",
		)
	}
	if c.req.Address != "" && c.req.Port != 0 {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and port for stage",
		)
	}
	c.inferStageAddress()
	return nil
}

func (c *CreateStageContext) inferStageAddress() {
	address := c.req.Address
	// If address is empty, fill it from req host and port.
	if address == "" {
		host, port := c.req.Host, c.req.Port
		if host == "" {
			host = defaultStageHost
		}
		if port == 0 {
			port = defaultStagePort
		}
		address = fmt.Sprintf("%s:%d", host, port)
	}
	c.inferredAddress = address
}

func (c *CreateStageContext) persist() error {
	var err error

	c.orchestration.Stages = append(c.orchestration.Stages, c.req.Name)
	if err = c.txnHelper.SaveOrchestration(c.orchestration); err != nil {
		return errdefs.PrependMsg(err, "save orchestration")
	}

	s := &api.Stage{
		Name:          c.req.Name,
		Phase:         api.StagePending,
		Service:       c.req.Service,
		Rpc:           c.req.Rpc,
		Address:       c.inferredAddress,
		Orchestration: c.orchestration.Name,
		Asset:         c.req.Asset,
	}
	if err = c.txnHelper.SaveStage(s); err != nil {
		return errdefs.PrependMsg(err, "save link")
	}
	return nil
}

type CreateLinkContext struct {
	// req contains the received request. This request should not be changed.
	req *api.CreateLinkRequest

	orchestration *api.Orchestration
	source        *api.Stage
	target        *api.Stage

	txnHelper *TxnHelper
}

func newCreateLinkContext(
	req *api.CreateLinkRequest,
	txnHelper *TxnHelper,
) *CreateLinkContext {
	return &CreateLinkContext{
		req:       req,
		txnHelper: txnHelper,
	}
}

func (c *CreateLinkContext) validateAndComplete() error {
	var orchestrationName api.OrchestrationName
	if ok, err := util.ArgNotNil(c.req, "req"); !ok {
		return err
	}
	if !IsValidLinkName(c.req.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", c.req.Name)
	}
	if c.txnHelper.ContainsLink(c.req.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"link '%v' already exists",
			c.req.Name,
		)
	}
	orchestrationName = c.req.Orchestration
	if orchestrationName == "" {
		orchestrationName = DefaultOrchestrationName
	}
	if !c.txnHelper.ContainsOrchestration(orchestrationName) {
		return errdefs.NotFoundWithMsg(
			"orchestration '%v' not found",
			orchestrationName,
		)
	}
	c.orchestration = &api.Orchestration{}
	err := c.txnHelper.LoadOrchestration(c.orchestration, orchestrationName)
	if err != nil {
		return errdefs.PrependMsg(err, "load orchestration")
	}
	if c.req.SourceStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty source stage name")
	}
	if c.req.TargetStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty target stage name")
	}
	if !c.txnHelper.ContainsStage(c.req.SourceStage) {
		return errdefs.NotFoundWithMsg(
			"source stage '%v' not found",
			c.req.SourceStage,
		)
	}
	c.source = &api.Stage{}
	if err := c.txnHelper.LoadStage(c.source, c.req.SourceStage); err != nil {
		return errdefs.PrependMsg(err, "load source stage")
	}
	if !c.txnHelper.ContainsStage(c.req.TargetStage) {
		return errdefs.NotFoundWithMsg(
			"target stage '%v' not found",
			c.req.TargetStage,
		)
	}
	c.target = &api.Stage{}
	if err := c.txnHelper.LoadStage(c.target, c.req.TargetStage); err != nil {
		return errdefs.PrependMsg(err, "load source stage")
	}
	if c.source.Orchestration != orchestrationName {
		return errdefs.FailedPreconditionWithMsg(
			"orchestration for link '%s' is '%s' but source stage is registered in '%s'.",
			c.req.Name,
			orchestrationName,
			c.source.Orchestration,
		)
	}
	if c.target.Orchestration != orchestrationName {
		return errdefs.FailedPreconditionWithMsg(
			"orchestration for link '%s' is '%s' but target stage is registered in '%s'.",
			c.req.Name,
			orchestrationName,
			c.target.Orchestration,
		)
	}
	if c.source.Phase != api.StagePending {
		return errdefs.FailedPreconditionWithMsg(
			"source stage is not in Pending phase for link '%s'.",
			c.req.Name,
		)
	}
	if c.target.Phase != api.StagePending {
		return errdefs.FailedPreconditionWithMsg(
			"target stage is not in Pending phase for link '%s'.",
			c.req.Name,
		)
	}
	return nil
}

func (c *CreateLinkContext) persist() error {
	var err error

	c.orchestration.Links = append(c.orchestration.Links, c.req.Name)
	if err = c.txnHelper.SaveOrchestration(c.orchestration); err != nil {
		return errdefs.PrependMsg(err, "save orchestration")
	}

	l := &api.Link{
		Name:          c.req.Name,
		SourceStage:   c.req.SourceStage,
		SourceField:   c.req.SourceField,
		TargetStage:   c.req.TargetStage,
		TargetField:   c.req.TargetField,
		Orchestration: c.orchestration.Name,
	}
	if err = c.txnHelper.SaveLink(l); err != nil {
		return errdefs.PrependMsg(err, "save link")
	}
	return nil
}
