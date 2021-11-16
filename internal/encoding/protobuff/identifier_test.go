package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

func TestMarshalID(t *testing.T) {
	rand, _ := identifier.Rand(5)
	tests := []struct {
		in identifier.Id
	}{
		{identifier.Empty()},
		{rand},
	}
	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res := MarshalID(in)
				assert.DeepEqual(t, res.Val, in.Val, "Id Val")
			})
	}
}

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

		t.Run(
			name, func(t *testing.T) {
				res := UnmarshalId(in)
				assert.DeepEqual(t, in.Val, res.Val, "Id Val")
			})
	}
}
