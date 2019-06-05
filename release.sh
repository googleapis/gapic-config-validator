#!/bin/bash
#
# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

VERSION=0.2.5

# linux-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/protoc-gen-gapic-validator
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/gapic-error-conformance
tar -czf gapic-config-validator-$VERSION-linux-amd64.tar.gz protoc-gen-gapic-validator gapic-error-conformance
rm -f gapic-error-conformance protoc-gen-gapic-validator

# linux-arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ./cmd/protoc-gen-gapic-validator
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ./cmd/gapic-error-conformance
tar -czf gapic-config-validator-$VERSION-linux-arm.tar.gz protoc-gen-gapic-validator gapic-error-conformance
rm -f gapic-error-conformance protoc-gen-gapic-validator

# darwin-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ./cmd/protoc-gen-gapic-validator
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ./cmd/gapic-error-conformance
tar -czf gapic-config-validator-$VERSION-darwin-amd64.tar.gz protoc-gen-gapic-validator gapic-error-conformance
rm -f gapic-error-conformance protoc-gen-gapic-validator

# windows-amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./cmd/protoc-gen-gapic-validator
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./cmd/gapic-error-conformance
tar -czf gapic-config-validator-$VERSION-windows-amd64.tar.gz protoc-gen-gapic-validator.exe gapic-error-conformance.exe
rm -f gapic-error-conformance.exe protoc-gen-gapic-validator.exe

# build & tag image
make image
docker tag \
  gcr.io/gapic-images/gapic-config-validator \
  gcr.io/gapic-images/gapic-config-validator:$VERSION
