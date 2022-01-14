package output

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/flow/flow"
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
	flow1, err := flow.NewFlow(link.New(linkName1, "", "", "", ""))
	assert.NilError(t, err, "create flow 1")
	flow2, err := flow.NewFlow(link.New(linkName2, "", "", "", ""))
	assert.NilError(t, err, "create flow 2")
	flow3, err := flow.NewFlow(link.New(linkName3, "", "", "", ""))
	assert.NilError(t, err, "create flow 3")

	cfg := NewOutputCfg()
	cfg.flows = append(cfg.flows, flow1)
	cfg.flows = append(cfg.flows, flow2)
	cfg.flows = append(cfg.flows, flow3)

	cfg.UnregisterIfExists(flow2)

	assert.Equal(t, 2, len(cfg.flows))
	assert.Equal(t, linkName1, cfg.flows[0].Link.Name(), "correct Link 1")
	assert.Equal(t, linkName3, cfg.flows[1].Link.Name(), "correct Link 3")
}

func TestOutputCfg_unregisterIfExists_doesNotExist(t *testing.T) {
	const (
		linkName1 apitypes.LinkName = "link-name-1"
		linkName2 apitypes.LinkName = "link-name-2"
		linkName3 apitypes.LinkName = "link-name-3"
		linkName4 apitypes.LinkName = "link-name-4"
	)
	flow1, err := flow.NewFlow(link.New(linkName1, "", "", "", ""))
	assert.NilError(t, err, "create flow 1")
	flow2, err := flow.NewFlow(link.New(linkName2, "", "", "", ""))
	assert.NilError(t, err, "create flow 2")
	flow3, err := flow.NewFlow(link.New(linkName3, "", "", "", ""))
	assert.NilError(t, err, "create flow 3")
	flow4, err := flow.NewFlow(link.New(linkName4, "", "", "", ""))
	assert.NilError(t, err, "create flow 4")

	cfg := NewOutputCfg()
	cfg.flows = append(cfg.flows, flow1)
	cfg.flows = append(cfg.flows, flow2)
	cfg.flows = append(cfg.flows, flow4)

	cfg.UnregisterIfExists(flow3)

	assert.Equal(t, 3, len(cfg.flows))
	assert.Equal(t, linkName1, cfg.flows[0].Link.Name(), "correct Link 1")
	assert.Equal(t, linkName2, cfg.flows[1].Link.Name(), "correct Link 2")
	assert.Equal(t, linkName4, cfg.flows[2].Link.Name(), "correct Link 4")
}
