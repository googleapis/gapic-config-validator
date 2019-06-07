// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/genproto/googleapis/rpc/status"
	
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-config-validator/internal/validator"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

var (
	plug    string
	opts    string
	verbose bool
)

func init() {
	flag.StringVar(&plug, "plugin", "", "path to the plugin binary to execute")
	flag.StringVar(&opts, "plugin_opts", "", "comma-delimited list of options for the plugin")
	flag.BoolVar(&verbose, "verbose", false, "log all response error")

	flag.Parse()

	if plug == "" {
		log.Fatalln("missing required flag -plugin")
	}
}

func main() {
	// load the testing scenarios containing CodeGeneratorRequests
	scenarios := scenarios()

	// run conformance evaluation
	var failed bool
	for _, s := range scenarios {
		// execute CodeGeneratorRequest with both plugins
		vResp, pResp, err := gen(s.req)
		if err != nil {
			log.Fatal(err)
		}
		if verbose {
			log.Println("validator:", vResp.GetError())
			log.Println("plugin:   ", pResp.GetError())
		}

		// validator & plugin response error messages
		if diff := compare(vResp.GetError(), pResp.GetError()); diff != nil {
			fmt.Println()
			fmt.Println(s.name, diff)
			failed = true
		}
	}

	if failed {
		os.Exit(1)
	}
}

// gen executes the CodeGeneratorRequest with the gapic-config-validator
// and the plugin named via the -plugin flag, and returns both responses.
func gen(req *plugin.CodeGeneratorRequest) (vResp, pResp *plugin.CodeGeneratorResponse, err error) {
	vResp, err = validator.Validate(req)
	if err != nil {
		log.Fatal(err)
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(plug)
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	_, err = in.Write(reqData)
	if err != nil {
		log.Fatal(err)
	}
	in.Close()

	resData, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	pResp = &plugin.CodeGeneratorResponse{}
	err = proto.Unmarshal(resData, pResp)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// compare checks if the plugin response exists in the validator response
// and produces an error with the diff in message content if not.
func compare(v, p string) error {
	if strings.Contains(v, p) {
		return nil
	}

	diff := lcsDiff([]string{p}, '+', strings.Split(v, ";"), '-')

	return fmt.Errorf("(+got,-want):\n%s", diff)
}

type scenario struct {
	name string
	req  *plugin.CodeGeneratorRequest
}

// scenarios builds the protos and CodeGeneratorRequest objects
// that define the error scenarios to test.
func scenarios() []scenario {
	// load config annotation common proto descriptors
	common := common()

	// missing default_host ServiceOption
	defHost := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, &descriptor.ServiceOptions{}, nil, nil, false)),
		Parameter:      proto.String(opts),
	}

	// missing response_type field in operation_info
	respInfoOpts := &descriptor.MethodOptions{}
	proto.SetExtension(
		respInfoOpts,
		longrunning.E_OperationInfo,
		&longrunning.OperationInfo{
			MetadataType: "Foo",
		},
	)

	missingRespInfo := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, respInfoOpts, nil, true)),
		Parameter:      proto.String(opts),
	}

	// missing MethodOption operation_info
	missingOpInfo := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, nil, nil, true)),
		Parameter:      proto.String(opts),
	}

	// unresolvable Message for response_type field in operation_info
	badRespInfoOpts := &descriptor.MethodOptions{}
	proto.SetExtension(
		respInfoOpts,
		longrunning.E_OperationInfo,
		&longrunning.OperationInfo{
			ResponseType: "Bad",
			MetadataType: "Foo",
		},
	)

	badRespInfo := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, badRespInfoOpts, nil, true)),
		Parameter:      proto.String(opts),
	}

	// unresolvable Message for metadata_type field in operation_info
	badMetaTypeOpts := &descriptor.MethodOptions{}
	proto.SetExtension(
		badMetaTypeOpts,
		longrunning.E_OperationInfo,
		&longrunning.OperationInfo{
			ResponseType: "Foo",
			MetadataType: "Bad",
		},
	)

	badMetaType := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, badMetaTypeOpts, nil, true)),
		Parameter:      proto.String(opts),
	}

	// unresolvable field in method_signature entry
	badSigFieldOpts := &descriptor.MethodOptions{}
	proto.SetExtension(
		badSigFieldOpts,
		annotations.E_MethodSignature,
		[]string{"a,bad"},
	)

	badSigField := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, badSigFieldOpts, nil, false)),
		Parameter:      proto.String(opts),
	}

	// unresolvable field in method_signature entry
	badCompSigOpts := &descriptor.MethodOptions{}
	proto.SetExtension(
		badCompSigOpts,
		annotations.E_MethodSignature,
		[]string{"a,bar.c"},
	)

	badCompSig := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, badCompSigOpts, nil, false)),
		Parameter:      proto.String(opts),
	}

	// unresolvable resource_reference
	refDNEOpts := &descriptor.FieldOptions{}
	proto.SetExtension(
		refDNEOpts,
		annotations.E_ResourceReference,
		&annotations.ResourceReference{
			Type: "DNE",
		},
	)

	refDNE := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, nil, refDNEOpts, false)),
		Parameter:      proto.String(opts),
	}

	// resource_reference leads to unannotated field
	unannotatedRefOpts := &descriptor.FieldOptions{}
	proto.SetExtension(
		unannotatedRefOpts,
		annotations.E_ResourceReference,
		&annotations.ResourceReference{
			Type: "Bar",
		},
	)

	unannotatedRef := &plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
		ProtoFile:      append(common, buildProto(nil, nil, nil, unannotatedRefOpts, false)),
		Parameter:      proto.String(opts),
	}

	return []scenario{
		{name: "missing default_host", req: defHost},
		{name: "missing LRO response_type", req: missingRespInfo},
		{name: "missing LRO operation_info", req: missingOpInfo},
		{name: "unresolvable LRO response_type", req: badRespInfo},
		{name: "unresolvable LRO metadata_type", req: badMetaType},
		{name: "bad method_signature field", req: badSigField},
		{name: "repeated nested field component in method_signature", req: badCompSig},
		{name: "unresolvable Message for resource_reference", req: refDNE},
		{name: "resource_reference to unannotated field", req: unannotatedRef},
	}
}

// buildProto constructs a foo.proto with the given options (or valid defaults)
// that is used for conformance testing scenarios.
func buildProto(fopts *descriptor.FileOptions, sopts *descriptor.ServiceOptions, mopts *descriptor.MethodOptions, fdopts *descriptor.FieldOptions, lro bool) *descriptor.FileDescriptorProto {
	if fopts == nil {
		fopts = &descriptor.FileOptions{}
	}
	fopts.GoPackage = proto.String("foo.com/bar/v1;bar")

	file := builder.NewFile("foo.proto").SetPackageName("foo").SetOptions(fopts)

	// when ServiceOptions is nil, set to default valid value
	if sopts == nil {
		sopts = &descriptor.ServiceOptions{}
		proto.SetExtension(sopts, annotations.E_DefaultHost, proto.String("api.foo.com"))
	}

	serv := builder.NewService("FooService").SetOptions(sopts)

	// other Message builder
	barBuilder := builder.NewMessage("Bar")
	barAFieldBuilder := builder.NewField("a", builder.FieldTypeString())
	barCFieldBuilder := builder.NewField("c", builder.FieldTypeString())
	barBuilder.AddField(barAFieldBuilder)
	barBuilder.AddField(barCFieldBuilder)

	// input & output Message builders
	var in, out *builder.RpcType
	inBuilder := builder.NewMessage("Foo")
	aFieldBuilder := builder.NewField("a", builder.FieldTypeString()).SetOptions(fdopts)
	bFieldBuilder := builder.NewField("b", builder.FieldTypeString())
	barFieldBuilder := builder.NewField("bar", builder.FieldTypeMessage(barBuilder)).SetRepeated()

	inBuilder.AddField(aFieldBuilder)
	inBuilder.AddField(bFieldBuilder)
	inBuilder.AddField(barFieldBuilder)

	// default output to same as input, LRO can change it
	out = builder.RpcTypeMessage(inBuilder, false)
	in = builder.RpcTypeMessage(inBuilder, false)

	if lro {
		lroDesc, _ := desc.LoadMessageDescriptorForMessage(&longrunning.Operation{})

		out = builder.RpcTypeImportedMessage(lroDesc, false)
	}

	mthd := builder.NewMethod("CreateFoo", in, out).SetOptions(mopts)

	serv.AddMethod(mthd)
	file.AddService(serv)
	file.AddMessage(inBuilder)
	file.AddMessage(barBuilder)
	f, _ := file.Build()

	return f.AsFileDescriptorProto()
}

// common loads the common descriptor dependencies for the config annotations.
func common() []*descriptor.FileDescriptorProto {
	protoDesc, err := desc.LoadMessageDescriptorForMessage(&descriptor.FileDescriptorProto{})
	if err != nil {
		log.Fatal(err)
	}

	annoDesc, err := desc.LoadFieldDescriptorForExtension(annotations.E_Http)
	if err != nil {
		log.Fatal(err)
	}

	httpDesc, err := desc.LoadMessageDescriptorForMessage(&annotations.Http{})
	if err != nil {
		log.Fatal(err)
	}

	anyDesc, err := desc.LoadMessageDescriptorForMessage(&any.Any{})
	if err != nil {
		log.Fatal(err)
	}

	emptyDesc, err := desc.LoadMessageDescriptorForMessage(&empty.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	statusDesc, err := desc.LoadMessageDescriptorForMessage(&status.Status{})
	if err != nil {
		log.Fatal(err)
	}

	clientDesc, err := desc.LoadFieldDescriptorForExtension(annotations.E_DefaultHost)
	if err != nil {
		log.Fatal(err)
	}

	resDesc, err := desc.LoadFieldDescriptorForExtension(annotations.E_Resource)
	if err != nil {
		log.Fatal(err)
	}

	fieldBehavDesc, err := desc.LoadFieldDescriptorForExtension(annotations.E_FieldBehavior)
	if err != nil {
		log.Fatal(err)
	}

	lroDesc, err := desc.LoadFieldDescriptorForExtension(longrunning.E_OperationInfo)
	if err != nil {
		log.Fatal(err)
	}

	durDesc, err := desc.LoadMessageDescriptorForMessage(&duration.Duration{})
	if err != nil {
		log.Fatal(err)
	}

	return []*descriptor.FileDescriptorProto{
		protoDesc.GetFile().AsFileDescriptorProto(),
		annoDesc.GetFile().AsFileDescriptorProto(),
		httpDesc.GetFile().AsFileDescriptorProto(),
		anyDesc.GetFile().AsFileDescriptorProto(),
		emptyDesc.GetFile().AsFileDescriptorProto(),
		statusDesc.GetFile().AsFileDescriptorProto(),
		clientDesc.GetFile().AsFileDescriptorProto(),
		resDesc.GetFile().AsFileDescriptorProto(),
		fieldBehavDesc.GetFile().AsFileDescriptorProto(),
		lroDesc.GetFile().AsFileDescriptorProto(),
		durDesc.GetFile().AsFileDescriptorProto(),
	}
}

// lcsDiff is copied from github.com/googleapis/gapic-generator-go/internal/txtdiff .
func lcsDiff(aLines []string, aSign rune, bLines []string, bSign rune) string {
	// Algorithm is described by https://en.wikipedia.org/wiki/Longest_common_subsequence_problem.

	// We require O(n^2) space to memoize LCS. This is not great, however
	// imagine we have 10,000-line baseline; 1e4^2 = 1e8 ints ~= 1e9 bytes = 1GB.
	// Most development computers have more memory than this and
	// our baselines are orders of magnitude smaller; we should we fine.

	// The article uses 1-based index and use index 0 to refer to the conceptual empty element.
	// Instead of dancing around the index, we just create the empty element.
	aLines = append([]string{""}, aLines...)
	bLines = append([]string{""}, bLines...)

	c := make([][]int, len(aLines))
	for i := range c {
		c[i] = make([]int, len(bLines))
	}
	for i := 1; i < len(aLines); i++ {
		for j := 1; j < len(bLines); j++ {
			if aLines[i] == bLines[j] {
				c[i][j] = c[i-1][j-1] + 1
			} else if c[i][j-1] < c[i-1][j] {
				c[i][j] = c[i-1][j]
			} else {
				c[i][j] = c[i][j-1]
			}
		}
	}

	// The article uses recursion. I think iteration is more clear.
	var diff []string
	var sign []rune

	i := len(aLines) - 1
	j := len(bLines) - 1
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && aLines[i] == bLines[j] {
			diff = append(diff, aLines[i])
			sign = append(sign, ' ')
			i--
			j--
		} else if j > 0 && (i == 0 || c[i][j-1] >= c[i-1][j]) {
			diff = append(diff, bLines[j])
			sign = append(sign, bSign)
			j--
		} else if i > 0 && (j == 0 || c[i][j-1] < c[i-1][j]) {
			diff = append(diff, aLines[i])
			sign = append(sign, aSign)
			i--
		}
	}

	var sb strings.Builder
	for i := len(diff) - 1; i >= 0; i-- {
		sb.WriteRune(sign[i])
		sb.WriteByte(' ')
		sb.WriteString(diff[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}
