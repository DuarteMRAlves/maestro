package grpcw

import (
	"testing"

	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRegisterFile(t *testing.T) {
	var emptyRegistry ProtoRegistry

	withEmpty := emptyRegistry.RegisterFile(emptypb.File_google_protobuf_empty_proto)
	verifyEmptyMessage(t, withEmpty)

	withAny := withEmpty.RegisterFile(anypb.File_google_protobuf_any_proto)
	verifyEmptyMessage(t, withAny)
	verifyAnyMessage(t, withAny)

	withCustom := withAny.RegisterFile(unit.File_registry_proto)
	verifyEmptyMessage(t, withCustom)
	verifyAnyMessage(t, withCustom)
	verifyInputMessage(t, withCustom)
	verifyOutputMessage(t, withCustom)
	verifyCustomService(t, withCustom)
}

func verifyEmptyMessage(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("google.protobuf.Empty")]
	if !ok {
		t.Fatal("empty message descriptor missing.")
	}
	msgDesc, ok := msgDescIface.(protoreflect.MessageDescriptor)
	if !ok {
		t.Fatal("empty message descriptor is not protoreflect.MessageDescritptor")
	}
	if diff := cmp.Diff(protoreflect.FullName("google.protobuf.Empty"), msgDesc.FullName()); diff != "" {
		t.Fatalf("mismatch empty message descriptor name:\n%s", diff)
	}
	if diff := cmp.Diff(0, msgDesc.Fields().Len()); diff != "" {
		t.Fatalf("mismatch on empty message descriptor number of fields:\n%s", diff)
	}
}

func verifyAnyMessage(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("google.protobuf.Any")]
	if !ok {
		t.Fatal("any message descriptor missing.")
	}
	msgDesc, ok := msgDescIface.(protoreflect.MessageDescriptor)
	if !ok {
		t.Fatal("any message descriptor is not protoreflect.MessageDescritptor")
	}
	if diff := cmp.Diff(protoreflect.FullName("google.protobuf.Any"), msgDesc.FullName()); diff != "" {
		t.Fatalf("mismatch any message descriptor name:\n%s", diff)
	}
	if diff := cmp.Diff(2, msgDesc.Fields().Len()); diff != "" {
		t.Fatalf("mismatch on any message descriptor number of fields:\n%s", diff)
	}
	if diff := cmp.Diff(
		protoreflect.FullName("google.protobuf.Any.type_url"),
		msgDesc.Fields().ByNumber(1).FullName(),
	); diff != "" {
		t.Fatalf("mismatch on any message field 1 descriptor:\n%s", diff)
	}
	if diff := cmp.Diff(
		protoreflect.FullName("google.protobuf.Any.value"),
		msgDesc.Fields().ByNumber(2).FullName(),
	); diff != "" {
		t.Fatalf("mismatch on any message field 2 descriptor:\n%s", diff)
	}
}

func verifyInputMessage(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("unit.InputMessage")]
	if !ok {
		t.Fatal("input message descriptor missing.")
	}
	msgDesc, ok := msgDescIface.(protoreflect.MessageDescriptor)
	if !ok {
		t.Fatal("input message descriptor is not protoreflect.MessageDescritptor")
	}
	if diff := cmp.Diff(protoreflect.FullName("unit.InputMessage"), msgDesc.FullName()); diff != "" {
		t.Fatalf("mismatch input message descriptor name:\n%s", diff)
	}
	if diff := cmp.Diff(0, msgDesc.Fields().Len()); diff != "" {
		t.Fatalf("mismatch on input message descriptor number of fields:\n%s", diff)
	}
}

func verifyOutputMessage(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("unit.OutputMessage")]
	if !ok {
		t.Fatal("output message descriptor missing.")
	}
	msgDesc, ok := msgDescIface.(protoreflect.MessageDescriptor)
	if !ok {
		t.Fatal("output message descriptor is not protoreflect.MessageDescritptor")
	}
	if diff := cmp.Diff(protoreflect.FullName("unit.OutputMessage"), msgDesc.FullName()); diff != "" {
		t.Fatalf("mismatch output message descriptor name:\n%s", diff)
	}
	if diff := cmp.Diff(0, msgDesc.Fields().Len()); diff != "" {
		t.Fatalf("mismatch on output message descriptor number of fields:\n%s", diff)
	}
}

func verifyCustomService(t *testing.T, r ProtoRegistry) {
	descIface, ok := r.descs[protoreflect.FullName("unit.CustomService")]
	if !ok {
		t.Fatal("custom service descriptor missing")
	}
	srvDesc, ok := descIface.(protoreflect.ServiceDescriptor)
	if !ok {
		t.Fatal("custom service descriptor is not protoreflect.ServiceDescriptor")
	}
	if diff := cmp.Diff(protoreflect.FullName("unit.CustomService"), srvDesc.FullName()); diff != "" {
		t.Fatalf("mismatch custom service descriptor name:\n%s", diff)
	}
	if diff := cmp.Diff(2, srvDesc.Methods().Len()); diff != "" {
		t.Fatalf("mismatch on custom service descriptor number of methods:\n%s", diff)
	}

	unary := srvDesc.Methods().ByName("Unary")
	if unary == nil {
		t.Fatal("unary method not found in custom service")
	}

	withEmpty := srvDesc.Methods().ByName("WithEmpty")
	if withEmpty == nil {
		t.Fatal("with empty method not fouond in custom service")
	}
}

func TestRegisterAndFindFile(t *testing.T) {
	var registry ProtoRegistry

	registry = registry.
		RegisterFile(emptypb.File_google_protobuf_empty_proto).
		RegisterFile(anypb.File_google_protobuf_any_proto).
		RegisterFile(unit.File_registry_proto)

	emptyFd, err := registry.FindFileByPath("google/protobuf/empty.proto")
	if err != nil {
		t.Fatalf("find empty file descriptor: %s", err)
	}
	if emptyFd != emptypb.File_google_protobuf_empty_proto {
		t.Fatalf("empty file descriptor not the same struct")
	}
	anyFd, err := registry.FindFileByPath("google/protobuf/any.proto")
	if err != nil {
		t.Fatalf("find any file descriptor: %s", err)
	}
	if anyFd != anypb.File_google_protobuf_any_proto {
		t.Fatalf("any file descriptor not the same struct")
	}
	customFd, err := registry.FindFileByPath("registry.proto")
	if err != nil {
		t.Fatalf("find custom file descriptor: %s", err)
	}
	if customFd != unit.File_registry_proto {
		t.Fatalf("custom file descriptor not the same struct")
	}
}
