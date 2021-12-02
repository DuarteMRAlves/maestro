package resources

import (
	"fmt"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	assetName1  = "assetName1"
	assetImage1 = "assetImage1"
	assetName2  = "assetName2"
	assetImage2 = "assetImage2"
)

var fileContent = fmt.Sprintf(
	`
---
kind: asset
spec:
  optional: %v
  image: %v
---
kind: asset
spec:
  optional: %v
  image: %v`,
	assetName1,
	assetImage1,
	assetName2,
	assetImage2)

var expected = []*Resource{
	{
		Kind: "asset",
		Spec: map[string]string{"optional": assetName1, "image": assetImage1},
	},
	{
		Kind: "asset",
		Spec: map[string]string{"optional": assetName2, "image": assetImage2},
	},
}

func TestUnmarshalAsset(t *testing.T) {
	data := []byte(fileContent)

	resources, err := ParseBytes(data)
	assert.NilError(t, err, "err is not nil: %v", err)
	assert.DeepEqual(t, expected, resources)
}
