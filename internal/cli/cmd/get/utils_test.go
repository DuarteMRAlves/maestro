package get

import (
	"context"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"google.golang.org/grpc"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

// assetForNum deterministically creates an asset resource with the given
// number.
func assetForNum(num int) *apitypes.Asset {
	return &apitypes.Asset{
		Name:  assetNameForNum(num),
		Image: assetImageForNum(num),
	}
}

// stageForNum deterministically creates a stage resource with the given number.
// The associated asset name is the one used in assetForNum.
func stageForNum(num int) *apitypes.Stage {
	return &apitypes.Stage{
		Name:    stageNameForNum(num),
		Asset:   assetNameForNum(num),
		Service: stageServiceForNum(num),
		Rpc:     stageRpcForNum(num),
		Address: stageAddressForNum(num),
	}
}

// linkForNum deterministically creates a link resource with the given number.
// The associated source stage name is the one used in stageForNum with the num
// argument. The associated target stage name is the one used in the stageForNum
// with the num+1 argument.
func linkForNum(num int) *apitypes.Link {
	return &apitypes.Link{
		Name:        linkNameForNum(num),
		SourceStage: linkSourceStageForNum(num),
		SourceField: linkSourceFieldForNum(num),
		TargetStage: linkTargetStageForNum(num),
		TargetField: linkTargetFieldForNum(num),
	}
}

// orchestrationForNum deterministically creates an orchestration resource with
// the given number.
func orchestrationForNum(num int) *apitypes.Orchestration {
	return &apitypes.Orchestration{Name: orchestrationNameForNum(num)}
}

// assetNameForNum deterministically creates an asset name for a given number.
func assetNameForNum(num int) string {
	return fmt.Sprintf("asset-%v", num)
}

// assetImageForNum deterministically creates an image for a given number.
func assetImageForNum(num int) string {
	name := assetNameForNum(num)
	return fmt.Sprintf("image-%v", name)
}

// stageNameForNum deterministically creates a stage name for a given number.
func stageNameForNum(num int) string {
	return fmt.Sprintf("stage-%v", num)
}

// stageServiceForNum deterministically creates a stage service for a given
// number.
func stageServiceForNum(num int) string {
	return fmt.Sprintf("service-%v", num)
}

// stageRpcForNum deterministically creates a stage rpc name for a given
// number.
func stageRpcForNum(num int) string {
	return fmt.Sprintf("rpc-%v", num)
}

// stageAddressForNum deterministically creates a stage address for a given
// number.
func stageAddressForNum(num int) string {
	return fmt.Sprintf("address-%v", num)
}

// linkNameForNum deterministically creates a link name for a given number.
func linkNameForNum(num int) string {
	return fmt.Sprintf("link-%v", num)
}

// linkSourceStageForNum deterministically creates a link source stage for a
// given number.
func linkSourceStageForNum(num int) string {
	return stageNameForNum(num)
}

// linkSourceFieldForNum deterministically creates a link source field for a
// given number.
func linkSourceFieldForNum(num int) string {
	return fmt.Sprintf("source-field-%v", num)
}

// linkTargetStageForNum deterministically creates a link target stage for a
// given number.
func linkTargetStageForNum(num int) string {
	return stageNameForNum(num + 1)
}

// linkTargetFieldForNum deterministically creates a link target field for a
// given number.
func linkTargetFieldForNum(num int) string {
	return fmt.Sprintf("target-field-%v", num)
}

// orchestrationNameForNum deterministically creates an orchestration name for a
// given number.
func orchestrationNameForNum(num int) string {
	return fmt.Sprintf("orchestration-%v", num)
}

// populateAssets creates the assets in the server, asserting any occurred
// errors.
func populateAssets(
	t *testing.T,
	assets []*apitypes.Asset,
	addr string,
) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	// Create asset before executing command
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, a := range assets {
		err := c.CreateAsset(ctx, a)
		assert.NilError(t, err, "populate with assets")
	}

	return nil
}

// populateStages creates the stages in the server, asserting any occurred
// errors.
func populateStages(
	t *testing.T,
	stages []*apitypes.Stage,
	addr string,
) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	// Create asset before executing command
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, s := range stages {
		err := c.CreateStage(ctx, s)
		assert.NilError(t, err, "populate with stages")
	}
	return nil
}

// populateLinks creates the links in the server, asserting any occurred errors.
func populateLinks(
	t *testing.T,
	links []*apitypes.Link,
	addr string,
) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	// Create asset before executing command
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, l := range links {
		err := c.CreateLink(ctx, l)
		assert.NilError(t, err, "populate with links")
	}
	return nil
}

// populateOrchestrations creates the orchestrations in the server, asserting
// any occurred errors.
func populateOrchestrations(
	t *testing.T,
	orchestrations []*apitypes.Orchestration,
	addr string,
) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	// Create asset before executing command
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, o := range orchestrations {
		err := c.CreateOrchestration(ctx, o)
		assert.NilError(t, err, "populate with orchestrations")
	}
	return nil
}
