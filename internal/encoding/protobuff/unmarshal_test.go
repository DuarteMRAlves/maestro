package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

func TestUnmarshalId(t *testing.T) {
	tests := []struct {
		in *pb.Id
	}{
		{&pb.Id{Val: ""}},
		{&pb.Id{Val: "SomeId"}},
	}
	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(name, func(t *testing.T) {
			res := UnmarshalId(in)
			assert.DeepEqual(t, in.Val, res.Val, "Id Val")
		})
	}
}

func TestUnmarshalAsset(t *testing.T) {
	tests := []struct {
		in *pb.Asset
	}{
		{&pb.Asset{Id: &pb.Id{Val: "Some String"}, Name: "Asset Name"}},
		{&pb.Asset{Id: &pb.Id{Val: ""}, Name: "Asset Name"}},
		{&pb.Asset{Name: "Asset Name"}},
		{&pb.Asset{Id: &pb.Id{Val: "Some String"}}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(name, func(t *testing.T) {
			res, err := UnmarshalAsset(in)
			assert.DeepEqual(t, nil, err, "Error")
			if in.Id != nil {
				assert.DeepEqual(t, in.Id.Val, res.Id.Val, "Asset Id")
			} else {
				assert.DeepEqual(t, identifier.Empty(), res.Id, "Asset Empty Id")
			}
			assert.DeepEqual(t, in.Name, res.Name, "Asset Name")
		})
	}
}
