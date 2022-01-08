package flow

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"gotest.tools/v3/assert"
	"testing"
)

func TestInputCfg_unregisterIfExists_exists(t *testing.T) {
	const (
		linkName1 apitypes.LinkName = "link-name-1"
		linkName2 apitypes.LinkName = "link-name-2"
		linkName3 apitypes.LinkName = "link-name-3"
	)
	link1 := link.New(linkName1, "", "", "", "")
	link2 := link.New(linkName2, "", "", "", "")
	link3 := link.New(linkName3, "", "", "", "")

	cfg := newInputCfg()
	cfg.links = append(cfg.links, link1)
	cfg.links = append(cfg.links, link2)
	cfg.links = append(cfg.links, link3)

	cfg.unregisterIfExists(link2.Clone())

	assert.Equal(t, 2, len(cfg.links))
	assert.Equal(t, linkName1, cfg.links[0].Name(), "correct link 1")
	assert.Equal(t, linkName3, cfg.links[1].Name(), "correct link 3")
}

func TestInputCfg_unregisterIfExists_doesNotExist(t *testing.T) {
	const (
		linkName1 apitypes.LinkName = "link-name-1"
		linkName2 apitypes.LinkName = "link-name-2"
		linkName3 apitypes.LinkName = "link-name-3"
		linkName4 apitypes.LinkName = "link-name-4"
	)
	link1 := link.New(linkName1, "", "", "", "")
	link2 := link.New(linkName2, "", "", "", "")
	link3 := link.New(linkName3, "", "", "", "")
	link4 := link.New(linkName4, "", "", "", "")

	cfg := newInputCfg()
	cfg.links = append(cfg.links, link1)
	cfg.links = append(cfg.links, link2)
	cfg.links = append(cfg.links, link4)

	cfg.unregisterIfExists(link3.Clone())

	assert.Equal(t, 3, len(cfg.links))
	assert.Equal(t, linkName1, cfg.links[0].Name(), "correct link 1")
	assert.Equal(t, linkName2, cfg.links[1].Name(), "correct link 2")
	assert.Equal(t, linkName4, cfg.links[2].Name(), "correct link 4")
}
