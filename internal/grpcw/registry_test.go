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
	verifyEmptyMessageMap(t, withEmpty)

	withAny := withEmpty.RegisterFile(anypb.File_google_protobuf_any_proto)
	verifyEmptyMessageMap(t, withAny)
	verifyAnyMessageMap(t, withAny)

	withCustom := withAny.RegisterFile(unit.File_registry_proto)
	verifyEmptyMessageMap(t, withCustom)
	verifyAnyMessageMap(t, withCustom)
	verifyInputMessageMap(t, withCustom)
	verifyOutputMessageMap(t, withCustom)
	verifyCustomServiceMap(t, withCustom)
}

func TestRegisterAndFindDescriptor(t *testing.T) {
	var registry ProtoRegistry

	registry = registry.
		RegisterFile(emptypb.File_google_protobuf_empty_proto).
		RegisterFile(anypb.File_google_protobuf_any_proto).
		RegisterFile(unit.File_registry_proto)

	verifyEmptyMessageFind(t, registry)
	verifyAnyMessageFind(t, registry)
	verifyInputMessageFind(t, registry)
	verifyOutputMessageFind(t, registry)
	verifyCustomServiceFind(t, registry)
	verifyUnaryMethodFind(t, registry)
	verifyWithEmptyMethodFind(t, registry)
}

func verifyEmptyMessageMap(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("google.protobuf.Empty")]
	if !ok {
		t.Fatal("empty message descriptor missing.")
	}
	verifyEmptyMessage(t, msgDescIface)
}

func verifyEmptyMessageFind(t *testing.T, r ProtoRegistry) {
	msgDescIface, err := r.FindDescriptorByName(protoreflect.FullName("google.protobuf.Empty"))
	if err != nil {
		t.Fatalf("find empty message descriptor: %v", err)
	}
	verifyEmptyMessage(t, msgDescIface)
}

func verifyEmptyMessage(t *testing.T, msgDescIface protoreflect.Descriptor) {
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

func verifyAnyMessageMap(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("google.protobuf.Any")]
	if !ok {
		t.Fatal("any message descriptor missing.")
	}
	verifyAnyMessage(t, msgDescIface)
}

func verifyAnyMessageFind(t *testing.T, r ProtoRegistry) {
	msgDescIface, err := r.FindDescriptorByName(protoreflect.FullName("google.protobuf.Any"))
	if err != nil {
		t.Fatalf("find any message descriptor: %v", err)
	}
	verifyAnyMessage(t, msgDescIface)
}

func verifyAnyMessage(t *testing.T, msgDescIface protoreflect.Descriptor) {
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

func verifyInputMessageMap(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("unit.InputMessage")]
	if !ok {
		t.Fatal("input message descriptor missing.")
	}
	verifyInputMessage(t, msgDescIface)
}

func verifyInputMessageFind(t *testing.T, r ProtoRegistry) {
	msgDescIface, err := r.FindDescriptorByName(protoreflect.FullName("unit.InputMessage"))
	if err != nil {
		t.Fatalf("find input message descriptor: %v", err)
	}
	verifyInputMessage(t, msgDescIface)
}

func verifyInputMessage(t *testing.T, msgDescIface protoreflect.Descriptor) {
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

func verifyOutputMessageMap(t *testing.T, r ProtoRegistry) {
	msgDescIface, ok := r.descs[protoreflect.FullName("unit.OutputMessage")]
	if !ok {
		t.Fatal("output message descriptor missing.")
	}
	verifyOutputMessage(t, msgDescIface)
}

func verifyOutputMessageFind(t *testing.T, r ProtoRegistry) {
	msgDescIface, err := r.FindDescriptorByName(protoreflect.FullName("unit.OutputMessage"))
	if err != nil {
		t.Fatalf("find output message descriptor: %v", err)
	}
	verifyOutputMessage(t, msgDescIface)
}

func verifyOutputMessage(t *testing.T, msgDescIface protoreflect.Descriptor) {
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

func verifyCustomServiceMap(t *testing.T, r ProtoRegistry) {
	descIface, ok := r.descs[protoreflect.FullName("unit.CustomService")]
	if !ok {
		t.Fatal("custom service descriptor missing")
	}
	verifyCustomService(t, descIface)
}

func verifyCustomServiceFind(t *testing.T, r ProtoRegistry) {
	descIface, err := r.FindDescriptorByName(protoreflect.FullName("unit.CustomService"))
	if err != nil {
		t.Fatalf("find custom service descriptor: %v", err)
	}
	verifyCustomService(t, descIface)
}

func verifyCustomService(t *testing.T, descIface protoreflect.Descriptor) {
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

	verifyUnaryMethod(t, srvDesc.Methods().ByName("Unary"))
	verifyWithEmptyMethod(t, srvDesc.Methods().ByName("WithEmpty"))
}

func verifyUnaryMethodFind(t *testing.T, r ProtoRegistry) {
	descIface, err := r.FindDescriptorByName(protoreflect.FullName("unit.CustomService.Unary"))
	if err != nil {
		t.Fatalf("find unary method descriptor: %v", err)
	}
	verifyUnaryMethod(t, descIface)
}

func verifyUnaryMethod(t *testing.T, descIface protoreflect.Descriptor) {
	methodDesc, ok := descIface.(protoreflect.MethodDescriptor)
	if !ok {
		t.Fatal("unary method descriptor is not protoreflect.MethodDescriptor")
	}
	if diff := cmp.Diff(protoreflect.FullName("unit.CustomService.Unary"), methodDesc.FullName()); diff != "" {
		t.Fatalf("mismatch unary method descriptor name:\n%s", diff)
	}
	verifyInputMessage(t, methodDesc.Input())
	verifyOutputMessage(t, methodDesc.Output())
}

func verifyWithEmptyMethodFind(t *testing.T, r ProtoRegistry) {
	descIface, err := r.FindDescriptorByName(protoreflect.FullName("unit.CustomService.WithEmpty"))
	if err != nil {
		t.Fatalf("find with empty method descriptor: %v", err)
	}
	verifyWithEmptyMethod(t, descIface)
}

func verifyWithEmptyMethod(t *testing.T, descIface protoreflect.Descriptor) {
	methodDesc, ok := descIface.(protoreflect.MethodDescriptor)
	if !ok {
		t.Fatal("with empty method descriptor is not protoreflect.MethodDescriptor")
	}
	if diff := cmp.Diff(protoreflect.FullName("unit.CustomService.WithEmpty"), methodDesc.FullName()); diff != "" {
		t.Fatalf("mismatch with empty method descriptor name:\n%s", diff)
	}
	verifyEmptyMessage(t, methodDesc.Input())
	verifyOutputMessage(t, methodDesc.Output())
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
