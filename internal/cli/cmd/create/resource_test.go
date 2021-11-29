package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/test"
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
  name: %v
  image: %v
---
kind: asset
spec:
  name: %v
  image: %v`,
	assetName1,
	assetImage1,
	assetName2,
	assetImage2)

var expected = []*Resource{
	{
		Kind: "asset",
		Spec: map[string]string{"name": assetName1, "image": assetImage1},
	},
	{
		Kind: "asset",
		Spec: map[string]string{"name": assetName2, "image": assetImage2},
	},
}

func TestUnmarshalAsset(t *testing.T) {
	data := []byte(fileContent)

	resources, err := UnmarshalResources(data)
	test.IsNil(t, err, "err is not nil: %v", err)
	test.DeepEqual(t, expected, resources, "resources not equal")
}
