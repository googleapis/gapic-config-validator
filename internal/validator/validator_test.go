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
	"testing"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/longrunning"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"

	"github.com/googleapis/gapic-config-validator/internal/validator/testdata"
)

func TestValidate(t *testing.T) {
	msg, err := desc.LoadMessageDescriptorForMessage(&testdata.Msg{})
	if err != nil {
		t.Error(err)
	}

	req := &plugin.CodeGeneratorRequest{
		ProtoFile:      []*descriptor.FileDescriptorProto{msg.GetFile().AsFileDescriptorProto()},
		FileToGenerate: []string{"basic_test.proto"},
	}

	res, err := Validate(req)
	if err != nil {
		t.Error(err)
	}

	if res.GetError() != "" {
		t.Errorf("Validate: received unexpected error(s) %s", res.GetError())
	}
}

func TestValidateFile(t *testing.T) {
	var v validator
	missingOpts := &descriptor.ServiceOptions{}
	missingServ := builder.NewService("missingService").SetOptions(missingOpts)

	missing, err := builder.NewFile("missing").AddService(missingServ).Build()
	if err != nil {
		t.Error(err)
	}

	for _, tst := range []struct {
		name, want string
		file       *desc.FileDescriptor
	}{
		{name: "missing default_host in Service", want: fmt.Sprintf(missingDefaultHost, missingServ.GetName()), file: missing},
	} {
		if err := v.validate(tst.file); err != nil {
			t.Error(err)
		}

		if actual := v.resp.GetError(); actual != tst.want {
			t.Errorf("%s: got(%s) want(%s)", tst.name, actual, tst.want)
		}

		// reset resp.Error field between tests
		v.resp.Error = nil
	}
}

func TestValidateService(t *testing.T) {
	var v validator

	missingOpts := &descriptor.ServiceOptions{}
	missing, err := builder.NewService("missing").SetOptions(missingOpts).Build()
	if err != nil {
		t.Error(err)
	}

	emptValueOpts := &descriptor.ServiceOptions{}
	if err := proto.SetExtension(emptValueOpts, annotations.E_DefaultHost, proto.String("")); err != nil {
		t.Error(err)
	}
	empty, err := builder.NewService("empty").SetOptions(emptValueOpts).Build()
	if err != nil {
		t.Error(err)
	}

	validOpts := &descriptor.ServiceOptions{}
	if err := proto.SetExtension(validOpts, annotations.E_DefaultHost, proto.String("foo.bar.com")); err != nil {
		t.Error(err)
	}
	valid, err := builder.NewService("valid").SetOptions(validOpts).Build()
	if err != nil {
		t.Error(err)
	}

	none, err := builder.NewService("none").Build()
	if err != nil {
		t.Error(err)
	}

	for _, tst := range []struct {
		name, want string
		serv       *desc.ServiceDescriptor
	}{
		{name: "missing default_host", want: fmt.Sprintf(missingDefaultHost, missing.GetFullyQualifiedName()), serv: missing},
		{name: "empty default_host value", want: fmt.Sprintf(emptyDefaultHost, empty.GetFullyQualifiedName()), serv: empty},
		{name: "valid default_host value", want: "", serv: valid},
		{name: "no ServiceOptions", want: fmt.Sprintf(missingDefaultHost, none.GetFullyQualifiedName()), serv: none},
	} {
		if err := v.validateService(tst.serv); err != nil {
			t.Error(err)
		}

		if actual := v.resp.GetError(); actual != tst.want {
			t.Errorf("%s: got(%s) want(%s)", tst.name, actual, tst.want)
		}

		// reset resp.Error field between tests
		v.resp.Error = nil
	}
}

func TestValidateMethod_LRO(t *testing.T) {
	var v validator

	// parents for Method builder resolution
	fooFile := builder.NewFile("foo").SetPackageName("foo").AddMessage(builder.NewMessage("ValidMessage"))
	barFile := builder.NewFile("bar").SetPackageName("bar").AddMessage(builder.NewMessage("ImportedMessage"))

	serv := builder.NewService("service")
	fooFile.AddService(serv)

	f, err := fooFile.Build()
	if err != nil {
		t.Error(err)
	}

	b, err := barFile.Build()
	if err != nil {
		t.Error(err)
	}
	v.files = map[string]*desc.FileDescriptor{"foo": f, "bar": b}

	lroDesc, err := desc.LoadMessageDescriptorForMessage(&longrunning.Operation{})
	if err != nil {
		t.Error(err)
	}
	lro := builder.RpcTypeImportedMessage(lroDesc, false)

	noneBuilder := builder.NewMethod("none", lro, lro)
	serv.AddMethod(noneBuilder)
	none, err := noneBuilder.Build()
	if err != nil {
		t.Error(err)
	}

	missingOpts := &descriptor.MethodOptions{}
	missingBuilder := builder.NewMethod("missing", lro, lro).SetOptions(missingOpts)
	serv.AddMethod(missingBuilder)
	missing, err := missingBuilder.Build()
	if err != nil {
		t.Error(err)
	}

	missingTypesInfoOpts := &descriptor.MethodOptions{}
	if err := proto.SetExtension(missingTypesInfoOpts, longrunning.E_OperationInfo, &longrunning.OperationInfo{}); err != nil {
		t.Error(err)
	}
	missingTypesBuilder := builder.NewMethod("missingTypes", lro, lro).SetOptions(missingTypesInfoOpts)
	serv.AddMethod(missingTypesBuilder)
	missingTypes, err := missingTypesBuilder.Build()
	if err != nil {
		t.Error(err)
	}

	uInfo := &longrunning.OperationInfo{
		ResponseType: "FooMessage",
		MetadataType: "bar.BazMessage",
	}
	unresolvableOpts := &descriptor.MethodOptions{}
	if err := proto.SetExtension(unresolvableOpts, longrunning.E_OperationInfo, uInfo); err != nil {
		t.Error(err)
	}
	unresolvableBuilder := builder.NewMethod("unresolvable", lro, lro).SetOptions(unresolvableOpts)
	serv.AddMethod(unresolvableBuilder)
	unresolvable, err := unresolvableBuilder.Build()
	if err != nil {
		t.Error(err)
	}

	validInfo := &longrunning.OperationInfo{
		ResponseType: "ValidMessage",
		MetadataType: "bar.ImportedMessage",
	}
	validOpts := &descriptor.MethodOptions{}
	if err := proto.SetExtension(validOpts, longrunning.E_OperationInfo, validInfo); err != nil {
		t.Error(err)
	}
	validBuilder := builder.NewMethod("valid", lro, lro).SetOptions(validOpts)
	serv.AddMethod(validBuilder)
	valid, err := validBuilder.Build()
	if err != nil {
		t.Error(err)
	}

	for _, tst := range []struct {
		name, want string
		mthd       *desc.MethodDescriptor
	}{
		{name: "no Method options", want: fmt.Sprintf(missingLROInfo, none.GetFullyQualifiedName()), mthd: none},
		{name: "missing operation_info", want: fmt.Sprintf(missingLROInfo, missing.GetFullyQualifiedName()), mthd: missing},
		{name: "missing response_type & metadata_type", want: fmt.Sprintf(missingLROResponseType+"; "+missingLROMetadataType, missingTypes.GetFullyQualifiedName(), missingTypes.GetFullyQualifiedName()), mthd: missingTypes},
		{name: "unresolvable response_type & metadata_type", want: fmt.Sprintf(unresolvableLROResponseType+"; "+unresolvableLROMetadataType, uInfo.GetResponseType(), unresolvable.GetFullyQualifiedName(), uInfo.GetMetadataType(), unresolvable.GetFullyQualifiedName()), mthd: unresolvable},
		{name: "valid LRO operation_info", want: "", mthd: valid},
	} {
		if err := v.validateMethod(tst.mthd); err != nil {
			t.Error(err)
		}

		if actual := v.resp.GetError(); actual != tst.want {
			t.Errorf("%s: got(%s) want(%s)", tst.name, actual, tst.want)
		}

		// reset resp.Error field between tests
		v.resp.Error = nil
	}
}

func TestValidateMethod_MethodSignature(t *testing.T) {
	var v validator

	fooFile := builder.NewFile("foo")
	serv := builder.NewService("service")
	fooFile.AddService(serv)

	fooDesc, err := desc.LoadMessageDescriptorForMessage(&testdata.Foo{})
	if err != nil {
		t.Error(err)
	}
	payload := builder.RpcTypeImportedMessage(fooDesc, false)

	sigs := []string{
		"a,req",         // invalid, required after optional
		"bar.baz.biz.d", // invalid, field component is repeated
		"dne",           // invalid, top-level field doesn't exist
		"bar.dne.c",     // invalid, nested field component doesn't exist
		"bar.dne",       // invalid, nested field doesn't exist
		"a,bar.b",       // valid w/nested
		"bar.baz.biz",   // valid, last component is repeated
	}

	opts := &descriptor.MethodOptions{}
	if err := proto.SetExtension(opts, annotations.E_MethodSignature, sigs); err != nil {
		t.Error(err)
	}
	methodBuilder := builder.NewMethod("SignatureAll", payload, payload).SetOptions(opts)
	serv.AddMethod(methodBuilder)
	method, err := methodBuilder.Build()
	if err != nil {
		t.Error(err)
	}

	for _, tst := range []struct {
		name, want string
		mthd       *desc.MethodDescriptor
	}{
		{
			name: "method_signature all",
			want: fmt.Sprintf(requiredAfterOptional+"; "+fieldComponentRepeated+"; "+fieldDNE+"; "+fieldDNE+"; "+fieldDNE,
				// requiredAfterOptional
				method.GetFullyQualifiedName(),
				sigs[0],
				"req",
				// fieldComponentRepeated
				method.GetFullyQualifiedName(),
				sigs[1],
				// fieldDNE
				"dne",
				method.GetFullyQualifiedName(),
				sigs[2],
				fooDesc.GetFullyQualifiedName(),
				// fieldDNE
				"bar.dne.c",
				method.GetFullyQualifiedName(),
				sigs[3],
				fooDesc.GetFullyQualifiedName(),
				// fieldDNE
				"bar.dne",
				method.GetFullyQualifiedName(),
				sigs[4],
				fooDesc.GetFullyQualifiedName(),
			),
			mthd: method,
		},
	} {
		if err := v.validateMethod(tst.mthd); err != nil {
			t.Error(err)
		}

		if actual := v.resp.GetError(); actual != tst.want {
			t.Errorf("%s: got(%s) want(%s)", tst.name, actual, tst.want)
		}

		// reset resp.Error field between tests
		v.resp.Error = nil
	}
}

func TestValidateMessage(t *testing.T) {
	var v validator

	barDesc, err := desc.LoadMessageDescriptorForMessage(&testdata.Bar{})
	if err != nil {
		t.Error(err)
	}

	bazDesc, err := desc.LoadMessageDescriptorForMessage(&testdata.Baz{})
	if err != nil {
		t.Error(err)
	}

	bizDesc, err := desc.LoadMessageDescriptorForMessage(&testdata.Biz{})
	if err != nil {
		t.Error(err)
	}

	v.files = map[string]*desc.FileDescriptor{"annotated_test.proto": barDesc.GetFile()}

	for _, tst := range []struct {
		name, want string
		msg        *desc.MessageDescriptor
	}{
		{name: "valid top-level reference", want: "", msg: barDesc},
		{name: "unresolvable top-lvl resource ref", want: fmt.Sprintf(resRefNotValidMessage, "annotated.Biz.d", "Buz"), msg: bizDesc},
		{name: "unresolvable top-lvl resource ref, empty", want: fmt.Sprintf(resRefNotValidMessage, "annotated.Baz.c", ""), msg: bazDesc},
	} {
		if err := v.validateMessage(tst.msg); err != nil {
			t.Error(err)
		}

		if actual := v.resp.GetError(); actual != tst.want {
			t.Errorf("%s: got(%s) want(%s)", tst.name, actual, tst.want)
		}

		// reset resp.Error field between tests
		v.resp.Error = nil
	}
}
