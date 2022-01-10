package flow

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"gotest.tools/v3/assert"
	"testing"
)

func TestOutputCfg_unregisterIfExists_exists(t *testing.T) {
	const (
		linkName1 apitypes.LinkName = "link-name-1"
		linkName2 apitypes.LinkName = "link-name-2"
		linkName3 apitypes.LinkName = "link-name-3"
	)
	flow1, err := newFlow(link.New(linkName1, "", "", "", ""))
	assert.NilError(t, err, "create flow 1")
	flow2, err := newFlow(link.New(linkName2, "", "", "", ""))
	assert.NilError(t, err, "create flow 2")
	flow3, err := newFlow(link.New(linkName3, "", "", "", ""))
	assert.NilError(t, err, "create flow 3")

	cfg := newOutputCfg()
	cfg.flows = append(cfg.flows, flow1)
	cfg.flows = append(cfg.flows, flow2)
	cfg.flows = append(cfg.flows, flow3)

	cfg.unregisterIfExists(flow2)

	assert.Equal(t, 2, len(cfg.flows))
	assert.Equal(t, linkName1, cfg.flows[0].link.Name(), "correct link 1")
	assert.Equal(t, linkName3, cfg.flows[1].link.Name(), "correct link 3")
}

func TestOutputCfg_unregisterIfExists_doesNotExist(t *testing.T) {
	const (
		linkName1 apitypes.LinkName = "link-name-1"
		linkName2 apitypes.LinkName = "link-name-2"
		linkName3 apitypes.LinkName = "link-name-3"
		linkName4 apitypes.LinkName = "link-name-4"
	)
	flow1, err := newFlow(link.New(linkName1, "", "", "", ""))
	assert.NilError(t, err, "create flow 1")
	flow2, err := newFlow(link.New(linkName2, "", "", "", ""))
	assert.NilError(t, err, "create flow 2")
	flow3, err := newFlow(link.New(linkName3, "", "", "", ""))
	assert.NilError(t, err, "create flow 3")
	flow4, err := newFlow(link.New(linkName4, "", "", "", ""))
	assert.NilError(t, err, "create flow 4")

	cfg := newOutputCfg()
	cfg.flows = append(cfg.flows, flow1)
	cfg.flows = append(cfg.flows, flow2)
	cfg.flows = append(cfg.flows, flow4)

	cfg.unregisterIfExists(flow3)

	assert.Equal(t, 3, len(cfg.flows))
	assert.Equal(t, linkName1, cfg.flows[0].link.Name(), "correct link 1")
	assert.Equal(t, linkName2, cfg.flows[1].link.Name(), "correct link 2")
	assert.Equal(t, linkName4, cfg.flows[2].link.Name(), "correct link 4")
}
