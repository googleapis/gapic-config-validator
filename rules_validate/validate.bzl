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

load("@com_google_api_codegen//rules_gapic:gapic.bzl", "proto_custom_library")

def gapic_config_validation(name, srcs, gapic_yaml = None, **kwargs):
  file_args = {}

  if gapic_yaml:
    file_args = {
      gapic_yaml: "gapic-yaml",
    }

  proto_custom_library(
    name = name,
    deps = srcs,
    plugin = Label("//cmd/protoc-gen-gapic-validator"),
    plugin_file_args = file_args,
    # there isn't any output, but this is required by the rule
    output_suffix = ".srcjar",
    output_type = "gapic-validator",
    **kwargs
  )
