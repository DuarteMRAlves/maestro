package output

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/execution/connection"
	"gotest.tools/v3/assert"
	"testing"
)

func TestOutputCfg_unregisterIfExists_exists(t *testing.T) {
	const (
		linkName1 api.LinkName = "link-name-1"
		linkName2 api.LinkName = "link-name-2"
		linkName3 api.LinkName = "link-name-3"
	)
	conn1, err := connection.New(
		&api.Link{
			Name:        linkName1,
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	)
	assert.NilError(t, err, "create connection 1")
	conn2, err := connection.New(
		&api.Link{
			Name:        linkName2,
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	)
	assert.NilError(t, err, "create connection 2")
	conn3, err := connection.New(
		&api.Link{
			Name:        linkName3,
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	)
	assert.NilError(t, err, "create connection 3")

	cfg := NewCfg()
	cfg.connections = append(cfg.connections, conn1)
	cfg.connections = append(cfg.connections, conn2)
	cfg.connections = append(cfg.connections, conn3)

	cfg.UnregisterIfExists(conn2)

	assert.Equal(t, 2, len(cfg.connections))
	assert.Equal(t, linkName1, cfg.connections[0].LinkName(), "correct Link 1")
	assert.Equal(t, linkName3, cfg.connections[1].LinkName(), "correct Link 3")
}

func TestOutputCfg_unregisterIfExists_doesNotExist(t *testing.T) {
	const (
		linkName1 api.LinkName = "link-name-1"
		linkName2 api.LinkName = "link-name-2"
		linkName3 api.LinkName = "link-name-3"
		linkName4 api.LinkName = "link-name-4"
	)
	conn1, err := connection.New(
		&api.Link{
			Name:        linkName1,
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	)
	assert.NilError(t, err, "create connection 1")
	conn2, err := connection.New(
		&api.Link{
			Name:        linkName2,
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	)
	assert.NilError(t, err, "create connection 2")
	conn3, err := connection.New(
		&api.Link{
			Name:        linkName3,
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	)
	assert.NilError(t, err, "create connection 3")
	conn4, err := connection.New(
		&api.Link{
			Name:        linkName4,
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	)
	assert.NilError(t, err, "create connection 4")

	cfg := NewCfg()
	cfg.connections = append(cfg.connections, conn1)
	cfg.connections = append(cfg.connections, conn2)
	cfg.connections = append(cfg.connections, conn4)

	cfg.UnregisterIfExists(conn3)

	assert.Equal(t, 3, len(cfg.connections))
	assert.Equal(t, linkName1, cfg.connections[0].LinkName(), "correct Link 1")
	assert.Equal(t, linkName2, cfg.connections[1].LinkName(), "correct Link 2")
	assert.Equal(t, linkName4, cfg.connections[2].LinkName(), "correct Link 4")
}
