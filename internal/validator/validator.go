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
	"strings"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/longrunning"
)

const (
	// default_host related errors
	missingDefaultHost = "service %v is missing option google.api.default_host"
	emptyDefaultHost   = "service %v google.api.default_host is empty"

	// LRO operation_info related errors
	missingLROInfo              = "rpc %v returns google.longrunning.Operation but is missing option google.longrunning.operation_info"
	missingLROResponseType      = "rpc %v has google.longrunning.operation_info but is missing option google.longrunning.operation_info.response_type"
	missingLROMetadataType      = "rpc %v has google.longrunning.operation_info but is missing option google.longrunning.operation_info.metadata_type"
	unresolvableLROResponseType = "unable to resolve google.longrunning.operation_info.response_type value %v in rpc %v"
	unresolvableLROMetadataType = "unable to resolve google.longrunning.operation_info.metadata_type value %v in rpc %v"

	// method_signature related errors
	fieldDNE               = "field %v listed in rpc %v method signature entry (%v) does not exist in %v"
	requiredAfterOptional  = "rpc %v method signature entry (%v) lists required field %v after an optional field"
	fieldComponentRepeated = "rpc %v method signature entry field %v cannot be a field within a repeated field"
)

// Validate ensures that the given input protos have valid
// GAPIC configuration annotations.
func Validate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var v validator
	var err error

	v.files, err = desc.CreateFileDescriptors(req.GetProtoFile())
	if err != nil {
		return &v.resp, err
	}

	for _, name := range req.GetFileToGenerate() {
		rich, ok := v.files[name]
		if !ok {
			return &v.resp, fmt.Errorf("FileToGenerate (%s) did not have a rich descriptor", name)
		}

		err = v.validate(rich)
		if err != nil {
			// this is an error encountered while attempting to validate
			// the file, not a misconfiguration in the proto
			return &v.resp, err
		}
	}

	return &v.resp, nil
}

type validator struct {
	resp  plugin.CodeGeneratorResponse
	files map[string]*desc.FileDescriptor
}

// validate executes GAPIC configuration validation on the given
// rich file descriptor.
func (v *validator) validate(file *desc.FileDescriptor) error {
	// validate Services
	for _, serv := range file.GetServices() {
		if err := v.validateService(serv); err != nil {
			return err
		}
	}

	// validate Messages

	return nil
}

// validateService checks the Service-level configuration annotations
// and validates each of its methods.
func (v *validator) validateService(serv *desc.ServiceDescriptor) error {
	// validate google.api.default_host
	if opts := serv.GetServiceOptions(); opts == nil {
		v.addError(missingDefaultHost, serv.GetFullyQualifiedName())
	} else if eHost, err := proto.GetExtension(opts, annotations.E_DefaultHost); err != nil {
		v.addError(missingDefaultHost, serv.GetFullyQualifiedName())
	} else if host := *eHost.(*string); host == "" {
		v.addError(emptyDefaultHost, serv.GetFullyQualifiedName())
	}

	// validate Methods
	for _, mthd := range serv.GetMethods() {
		if err := v.validateMethod(mthd); err != nil {
			return err
		}
	}

	return nil
}

// validateMethod checks the Method-level configuration annotations.
func (v *validator) validateMethod(method *desc.MethodDescriptor) error {
	mFQN := method.GetFullyQualifiedName()

	// validate google.longrunning.operation_info
	if out := method.GetOutputType(); out.GetFullyQualifiedName() == "google.longrunning.Operation" {
		if opts := method.GetMethodOptions(); opts == nil {
			v.addError(missingLROInfo, mFQN)
		} else if eLRO, err := proto.GetExtension(opts, longrunning.E_OperationInfo); err != nil {
			v.addError(missingLROInfo, mFQN)
		} else {
			lro := eLRO.(*longrunning.OperationInfo)

			if res := lro.GetResponseType(); res == "" {
				v.addError(missingLROResponseType, mFQN)
			} else if msg := v.resolveLROType(res, method.GetFile()); msg == nil {
				v.addError(unresolvableLROResponseType, res, mFQN)
			}

			if meta := lro.GetMetadataType(); meta == "" {
				v.addError(missingLROMetadataType, mFQN)
			} else if msg := v.resolveLROType(meta, method.GetFile()); msg == nil {
				v.addError(unresolvableLROMetadataType, meta, mFQN)
			}
		}
	}

	// validate google.api.method_signature
	if eSig, err := proto.GetExtension(method.GetMethodOptions(), annotations.E_MethodSignature); err == nil {
		sigs := eSig.([]string)
		input := method.GetInputType()

		// validate each method signature entry
		for _, sig := range sigs {
			// individual method signatures are a comma-delimited string of fields
			fields := strings.Split(sig, ",")
			seenOptional := false

			for _, field := range fields {
				var f *desc.FieldDescriptor

				// nested field
				if split := strings.Split(field, "."); len(split) > 1 {
					msg := input

					// validate each level of nested field
					for ndx, component := range split {
						if f = msg.FindFieldByName(component); f == nil {
							break
						} else if f.IsRepeated() && ndx < len(split)-1 {
							v.addError(fieldComponentRepeated, method.GetFullyQualifiedName(), field)
							break
						}

						msg = f.GetMessageType()
					}
				} else {
					// top-level field
					f = input.FindFieldByName(field)
				}

				// field doesn't exist
				if f == nil {
					v.addError(fieldDNE, field, method.GetFullyQualifiedName(), sig, input.GetFullyQualifiedName())
				} else if eBehavior, err := proto.GetExtension(f.GetFieldOptions(), annotations.E_FieldBehavior); err == nil {
					behaviors := eBehavior.([]annotations.FieldBehavior)

					// validate order of required & optional fields
					for _, b := range behaviors {
						if b == annotations.FieldBehavior_REQUIRED && seenOptional {
							v.addError(requiredAfterOptional, method.GetFullyQualifiedName(), sig, field)
							break
						} else if b == annotations.FieldBehavior_OPTIONAL {
							seenOptional = true
							break
						}
					}
				} else {
					// if no field_behavior is specified, the field is optional
					seenOptional = true
				}
			}
		}
	}

	return nil
}

// addError adds the given validation error to the plugin response
// error field. If the response error field already exists, the new error
// is concatenated with a semicolon.
func (v *validator) addError(err string, info ...interface{}) {
	if len(info) > 0 {
		err = fmt.Sprintf(err, info...)
	}

	if existing := v.resp.GetError(); existing != "" {
		err = fmt.Sprintf("%s; %s", existing, err)
	}

	v.resp.Error = proto.String(err)
}

// resolveLROType finds the MessageDescriptor for a fully qualified name
// of an operation_info.response_type or operation_info.metadata_type.
func (v *validator) resolveLROType(name string, file *desc.FileDescriptor) *desc.MessageDescriptor {
	// not a fully qualified name, make it one and check in parent file
	if !strings.Contains(name, ".") {
		msg := file.FindMessage(file.GetPackage() + "." + name)
		if msg != nil {
			return msg
		}
	}

	// this Message must be imported, check validator's file set
	// iterating of the entire set isn't ideal, but necessary
	// when searching for single message in all of the protos
	for _, f := range v.files {
		if msg := f.FindMessage(name); msg != nil {
			return msg
		}
	}

	return nil
}
