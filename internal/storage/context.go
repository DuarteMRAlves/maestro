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
		orchestrationName = "default"
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

func (c *CreateStageContext) stage() *api.Stage {
	return &api.Stage{
		Name:          c.req.Name,
		Phase:         api.StagePending,
		Service:       c.req.Service,
		Rpc:           c.req.Rpc,
		Address:       c.inferredAddress,
		Orchestration: c.orchestration.Name,
		Asset:         c.req.Asset,
	}
}
