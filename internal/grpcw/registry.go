package grpcw

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoRegistry struct {
	descs map[protoreflect.FullName]protoreflect.Descriptor
	files map[string]protoreflect.FileDescriptor
}

// Register file registers the top level definitions of a file. The behaviour
// for duplicate descriptors is to replace with the new descriptors.
// The function also associates the file with its path. If multiple files have
// the same path, the last one is also kept.
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

	newFiles := make(map[string]protoreflect.FileDescriptor, len(r.files)+1)
	for k, v := range r.files {
		newFiles[k] = v
	}
	newFiles[f.Path()] = f
	return ProtoRegistry{descs: newDescs, files: newFiles}
}

func (r ProtoRegistry) FindFileByPath(p string) (protoreflect.Descriptor, error) {
	f, ok := r.files[p]
	if !ok {
		return nil, &errPathNotFound{Path: p}
	}
	return f, nil
}

func (r ProtoRegistry) String() string {
	return fmt.Sprintf("ProtoRegistry{\n\tdescriptors: %s,\n\tfiles: %s\n}", r.descs, r.files)
}

type errPathNotFound struct {
	Path string
}

func (e *errPathNotFound) Error() string {
	return fmt.Sprintf("path not found: %q", e.Path)
}
