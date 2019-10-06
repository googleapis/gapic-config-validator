package validator

import (
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// resolveResRefMessage finds the MessageDescriptor of a
// resource_reference's given type. It attempts to
// resolve the type in the local file before consulting
// a all available files in the file set.
func (v *validator) resolveResRefMessage(typ string, file *desc.FileDescriptor) *desc.MessageDescriptor {
	if typ == "" {
		return nil
	}

	// check local file first
	if m := v.resolveResRefType(typ, file); m != nil {
		return m
	}

	// check the whole world for resources
	//
	// iterating over the entire file set of
	// services is not ideal, but the unified
	// resource design will go through some churn
	for _, f := range v.files {
		if m := v.resolveResRefType(typ, f); m != nil {
			return m
		}
	}

	return nil
}

// resolveResRefType checks every message in the file
// for one that is annotated with the resource type
func (v *validator) resolveResRefType(typ string, f *desc.FileDescriptor) *desc.MessageDescriptor {
	eResDef, err := ext(f.GetFileOptions(), annotations.E_ResourceDefinition)
	if err == nil {
		resDefs := eResDef.([]*annotations.ResourceDescriptor)
		for _, res := range resDefs {
			if typ != res.GetType() {
				continue
			}

			// resource_definitions are orphaned, no backing Message, fake one
			name := typ[strings.Index(typ, "/")+1:]
			field := builder.NewField("name", builder.FieldTypeString())
			m, _ := builder.NewMessage(name).AddField(field).Build()

			return m
		}
	}

	for _, m := range f.GetMessageTypes() {
		eRes, err := ext(m.GetMessageOptions(), annotations.E_Resource)
		if err != nil {
			continue
		}
		res := eRes.(*annotations.ResourceDescriptor)

		if typ == res.GetType() {
			return m
		}
	}

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
