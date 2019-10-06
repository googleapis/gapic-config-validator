// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/googleapis/gapic-config-validator/internal/config"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/longrunning"
)

const (
	// default_host related errors
	missingDefaultHost = "service %q is missing option google.api.default_host"
	emptyDefaultHost   = "service %q google.api.default_host is empty"

	// LRO operation_info related errors
	missingLROInfo              = "rpc %q returns google.longrunning.Operation but is missing option google.longrunning.operation_info"
	missingLROResponseType      = "rpc %q has google.longrunning.operation_info but is missing option google.longrunning.operation_info.response_type"
	missingLROMetadataType      = "rpc %q has google.longrunning.operation_info but is missing option google.longrunning.operation_info.metadata_type"
	unresolvableLROResponseType = "unable to resolve google.longrunning.operation_info.response_type value %q in rpc %q"
	unresolvableLROMetadataType = "unable to resolve google.longrunning.operation_info.metadata_type value %q in rpc %q"

	// method_signature related errors
	fieldDNE               = "field %q listed in rpc %q method signature entry (%q) does not exist in %q"
	requiredAfterOptional  = "rpc %q method signature entry (%q) lists required field %q after an optional field"
	fieldComponentRepeated = "rpc %q method signature entry field %q cannot be a field within a repeated field"

	// resource reslated errors
	resRefNotValidResource  = "unable to resolve resource reference for field %q: value %q is not a valid resource"
	resRefFieldDNE          = "unable to resolve resource reference for field %q: field does not exist or is not defined on message %q"
	resRefInvalidTypeFormat = "resource_reference.(child_)type for field %q must be {service_name}/{resource_type_kind}"
	resMissingType          = "resource for message %q missing field google.api.resource.type"
	resInvalidTypeFormat    = "resource.(child_)type for message %q must be {service_name}/{resource_type_kind}"
	resTypeKindInvalid      = "resource_type_kind %q has invalid format, must match regexp [A-Z][a-zA-Z0-9]+"
	resTypeKindTooLong      = "resource_type_kind in message %q must not be longer than %d characters"
	resMissingPattern       = "field %q resource missing pattern definition"
	resMissingNameField     = "resource message %q missing a name field"

	maxCharRescTypeKind = 100
)

var (
	resourceTypeKindRegexp *regexp.Regexp
	wellKnownTypes         = map[string]bool{
		"cloudresourcemanager.googleapis.com/Project":      true,
		"cloudresourcemanager.googleapis.com/Organization": true,
		"cloudresourcemanager.googleapis.com/Folder":       true,
		"locations.googleapis.com/Location":                true,
		"cloudbilling.googleapis.com/BillingAccount":       true,
	}
)

// Validate ensures that the given input protos have valid
// GAPIC configuration annotations.
func Validate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var v validator
	var err error

	resourceTypeKindRegexp = regexp.MustCompile("[A-Z][a-zA-Z0-9]+")

	v.files, err = desc.CreateFileDescriptors(req.GetProtoFile())
	if err != nil {
		return &v.resp, err
	}

	err = v.parseParameters(req.GetParameter())
	if err != nil {
		return &v.resp, err
	}

	if v.gapic != nil {
		v.compare()
	}

	for _, name := range req.GetFileToGenerate() {
		rich, ok := v.files[name]
		if !ok {
			return &v.resp, fmt.Errorf("FileToGenerate (%s) did not have a rich descriptor", name)
		}

		v.validate(rich)
	}

	return &v.resp, nil
}

type validator struct {
	resp  plugin.CodeGeneratorResponse
	files map[string]*desc.FileDescriptor
	gapic *config.ConfigProto
}

// validate executes GAPIC configuration validation on the given
// rich file descriptor.
func (v *validator) validate(file *desc.FileDescriptor) {
	opts := file.GetFileOptions()
	eResDef, err := ext(opts, annotations.E_ResourceDefinition)
	if err == nil {
		resDefs := eResDef.([]*annotations.ResourceDescriptor)
		for _, res := range resDefs {
			v.validateResourceDescriptor(res, res.GetType())
		}
	}

	// validate Services
	for _, serv := range file.GetServices() {
		v.validateService(serv)
	}

	// validate Messages
	for _, msg := range file.GetMessageTypes() {
		v.validateMessage(msg)
	}
}

// validateService checks the Service-level configuration annotations
// and validates each of its methods.
func (v *validator) validateService(serv *desc.ServiceDescriptor) {
	// validate google.api.default_host
	if opts := serv.GetServiceOptions(); opts == nil {
		v.addError(missingDefaultHost, serv.GetFullyQualifiedName())
	} else if eHost, err := ext(opts, annotations.E_DefaultHost); err != nil {
		v.addError(missingDefaultHost, serv.GetFullyQualifiedName())
	} else if host := *eHost.(*string); host == "" {
		v.addError(emptyDefaultHost, serv.GetFullyQualifiedName())
	}

	// validate Methods
	for _, mthd := range serv.GetMethods() {
		v.validateMethod(mthd)
	}
}

// validateMethod checks the Method-level configuration annotations.
func (v *validator) validateMethod(method *desc.MethodDescriptor) {
	mFQN := method.GetFullyQualifiedName()

	// validate google.longrunning.operation_info
	if method.GetOutputType().GetFullyQualifiedName() == "google.longrunning.Operation" {
		if opts := method.GetMethodOptions(); opts == nil {
			v.addError(missingLROInfo, mFQN)
		} else if eLRO, err := ext(opts, longrunning.E_OperationInfo); err != nil {
			v.addError(missingLROInfo, mFQN)
		} else {
			lro := eLRO.(*longrunning.OperationInfo)

			if res := lro.GetResponseType(); res == "" {
				v.addError(missingLROResponseType, mFQN)
			} else if v.resolveMsgReference(res, method.GetFile()) == nil {
				v.addError(unresolvableLROResponseType, res, mFQN)
			}

			if meta := lro.GetMetadataType(); meta == "" {
				v.addError(missingLROMetadataType, mFQN)
			} else if v.resolveMsgReference(meta, method.GetFile()) == nil {
				v.addError(unresolvableLROMetadataType, meta, mFQN)
			}
		}
	}

	// validate google.api.method_signature
	if eSig, err := ext(method.GetMethodOptions(), annotations.E_MethodSignature); err == nil {
		sigs := eSig.([]string)
		input := method.GetInputType()

		// validate each method signature entry
		for _, sig := range sigs {
			// individual method signatures are a comma-delimited string of fields
			fields := strings.Split(sig, ",")

			for _, field := range fields {
				f := input.FindFieldByName(field)

				// nested field
				if split := strings.Split(field, "."); len(split) > 1 {
					msg := input

					// validate each level of nested field
					for ndx, component := range split {
						if f = msg.FindFieldByName(component); f == nil {
							break
						}

						if f.IsRepeated() && ndx < len(split)-1 {
							v.addError(
								fieldComponentRepeated,
								method.GetFullyQualifiedName(),
								field,
							)

							break
						}

						msg = f.GetMessageType()
					}
				}

				// field doesn't exist
				if f == nil {
					v.addError(
						fieldDNE,
						field,
						method.GetFullyQualifiedName(),
						sig,
						input.GetFullyQualifiedName(),
					)
				}
			}
		}
	}
}

func (v *validator) validateMessage(msg *desc.MessageDescriptor) {
	// validate message resource
	if eRes, err := ext(msg.GetMessageOptions(), annotations.E_Resource); err == nil {
		res := eRes.(*annotations.ResourceDescriptor)

		v.validateResourceDescriptor(res, msg.GetFullyQualifiedName())

		fname := "name"
		if n := res.GetNameField(); n != "" {
			fname = n
		}

		if f := msg.FindFieldByName(fname); f == nil {
			// missing resource name field
			v.addError(resMissingNameField, msg.GetFullyQualifiedName())
		}
	}

	for _, field := range msg.GetFields() {
		// validate individual resource reference
		if eRef, err := ext(field.GetFieldOptions(), annotations.E_ResourceReference); err == nil {
			v.validateResRef(eRef.(*annotations.ResourceReference), field)
		}
	}
}

// validateResourceDescriptor validates the resource_type_kind and pattern
// presence of a given ResourceDescriptor for the owner with the
// fully-qualified name fqn.
func (v *validator) validateResourceDescriptor(res *annotations.ResourceDescriptor, fqn string) {
	// missing resource.pattern
	if len(res.GetPattern()) == 0 {
		v.addError(resMissingPattern, fqn)
	}

	// missing resource.type
	typ := res.GetType()
	if typ == "" {
		v.addError(resMissingType, fqn)
		return
	}

	// validate resource.type format
	split := strings.Split(typ, "/")
	if len(split) != 2 {
		v.addError(resInvalidTypeFormat, fqn)
		return
	}

	v.validateRescTypeKind(split[1], fqn)
}

// validateRescTypeKind ensures that the resource_type_kind component
// of a resource.type conforms to the required format and length.
func (v *validator) validateRescTypeKind(rtk, fqn string) {
	if !resourceTypeKindRegexp.MatchString(rtk) {
		v.addError(resTypeKindInvalid, rtk)
	}

	if len(rtk) > maxCharRescTypeKind {
		v.addError(resTypeKindTooLong, fqn, maxCharRescTypeKind)
	}
}

// validateResRef ensures that the given resource_reference is resolvable
// within the field's file or the file set.
func (v *validator) validateResRef(ref *annotations.ResourceReference, field *desc.FieldDescriptor) {
	typ := ref.GetType()

	if typ == "" {
		typ = ref.GetChildType()
	}

	// check well-known types
	if wellKnownTypes[typ] || typ == "*" {
		return
	}

	if split := strings.Split(typ, "/"); len(split) != 2 {
		v.addError(resRefInvalidTypeFormat, field.GetFullyQualifiedName())
		return
	}

	refMsg := v.resolveResRefMessage(typ, field.GetFile())

	if refMsg == nil {
		v.addError(resRefNotValidResource, field.GetFullyQualifiedName(), typ)
	}
}

// addError adds the given validation error to the plugin response
// error field. If the response error field already exists, the new error
// is concatenated with a semicolon.
func (v *validator) addError(err string, info ...interface{}) {
	if len(info) > 0 {
		err = fmt.Sprintf(err, info...)
	}

	err = fmt.Sprintf("%s\n%s", v.resp.GetError(), err)

	v.resp.Error = proto.String(err)
}

// ext wraps proto.GetExtension
func ext(pb proto.Message, eDesc *proto.ExtensionDesc) (interface{}, error) {
	return proto.GetExtension(pb, eDesc)
}
