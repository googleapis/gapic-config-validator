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

package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-config-validator/internal/validator"
)

func main() {
	reqBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var req plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(reqBytes, &req); err != nil {
		log.Fatal(err)
	}

	resp, err := validator.Validate(&req)
	if err != nil {
		resp.Error = proto.String(err.Error())
	}

	outBytes, err := proto.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(outBytes); err != nil {
		log.Fatal(err)
	}
}
