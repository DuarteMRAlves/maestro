package grpcw

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoRegistry struct {
	descs map[protoreflect.FullName]protoreflect.Descriptor
}

// Register file registers the top level definitions of a file. The behaviour
// for duplicate descriptors is to replace with the new descriptors.
func (r ProtoRegistry) RegisterFile(f protoreflect.FileDescriptor) ProtoRegistry {
	var toAdd []protoreflect.Descriptor

	enums := f.Enums()
	for i := 0; i < enums.Len(); i++ {
		enum := enums.Get(i)
		toAdd = append(toAdd, enum)
		values := enum.Values()
		for j := 0; j < values.Len(); j++ {
			toAdd = append(toAdd, values.Get(i))
		}
	}

	msgs := f.Messages()
	for i := 0; i < msgs.Len(); i++ {
		toAdd = append(toAdd, msgs.Get(i))
	}

	exts := f.Extensions()
	for i := 0; i < exts.Len(); i++ {
		toAdd = append(toAdd, exts.Get(i))
	}

	srvs := f.Services()
	for i := 0; i < srvs.Len(); i++ {
		toAdd = append(toAdd, srvs.Get(i))
	}

	newDescs := make(map[protoreflect.FullName]protoreflect.Descriptor, len(r.descs)+len(toAdd))
	for k, v := range r.descs {
		newDescs[k] = v
	}
	for _, d := range toAdd {
		newDescs[d.FullName()] = d
	}
	return ProtoRegistry{descs: newDescs}
}

func (p ProtoRegistry) String() string {
	return fmt.Sprintf("ProtoRegistry(%s)", p.descs)
}
