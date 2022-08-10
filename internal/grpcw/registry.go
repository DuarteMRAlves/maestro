package grpcw

import (
	"fmt"
	"strings"

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

func (r ProtoRegistry) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) {
	var (
		desc        protoreflect.Descriptor
		prefix      protoreflect.FullName
		suffixParts []protoreflect.Name
	)

	// Start by full name and incrementally search for a smaller prefix
	for prefix = name; prefix != ""; prefix = prefix.Parent() {
		var ok bool
		if desc, ok = r.descs[prefix]; ok {
			break
		}
	}

	if prefix == "" {
		return nil, &errNameNotFound{Name: name}
	}

	// Did not match full name
	if prefix != name {
		suffixParts = fullNameParts(name[len(prefix)+len("."):])
	}

	switch desc := desc.(type) {
	case protoreflect.EnumDescriptor:
		if desc.FullName() == name {
			return desc, nil
		}
	case protoreflect.EnumValueDescriptor:
		if desc.FullName() == name {
			return desc, nil
		}
	case protoreflect.MessageDescriptor:
		if desc.FullName() == name {
			return desc, nil
		}
		if desc := findDescriptorInMessage(desc, suffixParts...); desc != nil && desc.FullName() == name {
			return desc, nil
		}
	case protoreflect.ExtensionDescriptor:
		if desc.FullName() == name {
			return desc, nil
		}
	case protoreflect.ServiceDescriptor:
		if desc.FullName() == name {
			return desc, nil
		}
		if d := desc.Methods().ByName(suffixParts[0]); d != nil && d.FullName() == name {
			return d, nil
		}
	}
	return nil, &errNameNotFound{Name: name}
}

func fullNameParts(n protoreflect.FullName) []protoreflect.Name {
	parts := strings.Split(string(n), ".")
	names := make([]protoreflect.Name, 0, len(parts))
	for _, p := range parts {
		names = append(names, protoreflect.Name(p))
	}
	return names
}

func findDescriptorInMessage(md protoreflect.MessageDescriptor, nameParts ...protoreflect.Name) protoreflect.Descriptor {
	if len(nameParts) == 0 {
		return nil
	}
	name := nameParts[0]
	if len(nameParts) == 1 {
		if ed := md.Enums().ByName(name); ed != nil {
			return ed
		}
		for i := 0; i < md.Enums().Len(); i++ {
			if vd := md.Enums().Get(i).Values().ByName(name); vd != nil {
				return vd
			}
		}
		if xd := md.Extensions().ByName(name); xd != nil {
			return xd
		}
		if fd := md.Fields().ByName(name); fd != nil {
			return fd
		}
		if od := md.Oneofs().ByName(name); od != nil {
			return od
		}
	}
	if md := md.Messages().ByName(name); md != nil {
		if len(nameParts) == 1 {
			return md
		}
		return findDescriptorInMessage(md, nameParts[1:]...)
	}
	return nil
}

func (r ProtoRegistry) FindFileByPath(p string) (protoreflect.FileDescriptor, error) {
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

type errNameNotFound struct {
	Name protoreflect.FullName
}

func (e *errNameNotFound) Error() string {
	return fmt.Sprintf("name not found: %q", e.Name)
}
