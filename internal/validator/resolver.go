package validator

import (
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// resolveResRefMessage finds the MessageDescriptor of a
// resource_reference's given type. It attempts to
// resolve the type in the local file before consulting
// a list of known resource types.
func (v *validator) resolveResRefMessage(typ, serv string, file *desc.FileDescriptor) *desc.MessageDescriptor {
	if typ == "" {
		return nil
	}

	// check local file
	if m := file.FindMessage(file.GetPackage() + "." + typ); m != nil {
		return m
	}

	// full := serv + "/" + typ

	// check configured types using serv
	//
	// TODO(ndietz)

	return nil
}

// resolveMsgReference finds the MessageDescriptor for a fully qualified name
// of an operation_info.response_type or operation_info.metadata_type.
func (v *validator) resolveMsgReference(name string, file *desc.FileDescriptor) *desc.MessageDescriptor {
	if name == "" {
		return nil
	}

	// not a fully qualified name, make it one and check in parent file
	//
	// TODO(ndietz) this will break if the name refs a nested message
	// in the parent file
	if !strings.Contains(name, ".") {
		if msg := file.FindMessage(file.GetPackage() + "." + name); msg != nil {
			return msg
		}
	}

	// this Message must be imported, check validator's file set.
	// Iterating over the entire set isn't ideal, but necessary
	// when searching for single message name in all protos
	for _, f := range v.files {
		if msg := f.FindMessage(name); msg != nil {
			return msg
		}
	}

	return nil
}
