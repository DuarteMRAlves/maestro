package input

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/flow/connection"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"gotest.tools/v3/assert"
	"testing"
)

func TestInputCfg_unregisterIfExists_exists(t *testing.T) {
	const (
		linkName1 apitypes.LinkName = "link-name-1"
		linkName2 apitypes.LinkName = "link-name-2"
		linkName3 apitypes.LinkName = "link-name-3"
	)
	conn1, err := connection.New(orchestration.NewLink(linkName1, "", "", "", ""))
	assert.NilError(t, err, "create connection 1")
	conn2, err := connection.New(orchestration.NewLink(linkName2, "", "", "", ""))
	assert.NilError(t, err, "create connection 2")
	conn3, err := connection.New(orchestration.NewLink(linkName3, "", "", "", ""))
	assert.NilError(t, err, "create connection 3")

	cfg := NewCfg()
	cfg.connections = append(cfg.connections, conn1)
	cfg.connections = append(cfg.connections, conn2)
	cfg.connections = append(cfg.connections, conn3)

	cfg.UnregisterIfExists(conn2)

	assert.Equal(t, 2, len(cfg.connections))
	assert.Equal(t, linkName1, cfg.connections[0].LinkName(), "correct link 1")
	assert.Equal(t, linkName3, cfg.connections[1].LinkName(), "correct link 3")
}

func TestInputCfg_unregisterIfExists_doesNotExist(t *testing.T) {
	const (
		linkName1 apitypes.LinkName = "link-name-1"
		linkName2 apitypes.LinkName = "link-name-2"
		linkName3 apitypes.LinkName = "link-name-3"
		linkName4 apitypes.LinkName = "link-name-4"
	)
	conn1, err := connection.New(orchestration.NewLink(linkName1, "", "", "", ""))
	assert.NilError(t, err, "create connection 1")
	conn2, err := connection.New(orchestration.NewLink(linkName2, "", "", "", ""))
	assert.NilError(t, err, "create connection 2")
	conn3, err := connection.New(orchestration.NewLink(linkName3, "", "", "", ""))
	assert.NilError(t, err, "create connection 3")
	conn4, err := connection.New(orchestration.NewLink(linkName4, "", "", "", ""))
	assert.NilError(t, err, "create connection 4")

	cfg := NewCfg()
	cfg.connections = append(cfg.connections, conn1)
	cfg.connections = append(cfg.connections, conn2)
	cfg.connections = append(cfg.connections, conn4)

	cfg.UnregisterIfExists(conn3)

	assert.Equal(t, 3, len(cfg.connections))
	assert.Equal(t, linkName1, cfg.connections[0].LinkName(), "correct link 1")
	assert.Equal(t, linkName2, cfg.connections[1].LinkName(), "correct link 2")
	assert.Equal(t, linkName4, cfg.connections[2].LinkName(), "correct link 4")
}
