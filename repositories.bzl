# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def com_googleapis_gapic_config_validator_repositories():
    _maybe(
        http_archive,
        name = "com_google_protobuf",
        strip_prefix = "protobuf-3.12.0-rc2",
        urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.12.0-rc2.tar.gz"],
        sha256 = "afaa4f65e7e97adb10b32b7c699b7b6be4090912b471028ef0f40ccfb271f96a",
    )

    _maybe(
        http_archive,
        name = "com_google_api_codegen",
        strip_prefix = "gapic-generator-213545efb83861e1b41d2b7095973d3c8239cc74",
        urls = ["https://github.com/googleapis/gapic-generator/archive/213545efb83861e1b41d2b7095973d3c8239cc74.zip"],
    )

    _maybe(
        http_archive,
        name = "io_bazel_rules_go",
        urls = [
            "https://github.com/bazelbuild/rules_go/archive/v0.23.0.zip",
        ],
        strip_prefix = "rules_go-0.23.0",
        sha256 = "4707e6ba7c01fcfc4f0d340d123bc16e43c2b8ea3f307663d95712b36d2a0e88",
    )

    _maybe(
        http_archive,
        name = "bazel_gazelle",
        urls = [
            "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.0/bazel-gazelle-v0.21.0.tar.gz",
            "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.0/bazel-gazelle-v0.21.0.tar.gz",
        ],
        sha256 = "bfd86b3cbe855d6c16c6fce60d76bd51f5c8dbc9cfcaef7a2bb5c1aafd0710e8",
    )

def _maybe(repo_rule, name, strip_repo_prefix = "", **kwargs):
    if not name.startswith(strip_repo_prefix):
        return
    repo_name = name[len(strip_repo_prefix):]
    if repo_name in native.existing_rules():
        return
    repo_rule(name = repo_name, **kwargs)
