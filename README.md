gapic-config-validator
======================

A configuration validator, and associated plugin error conformance utility, for GAPIC config proto annotations.

The `protoc-gen-gapic-validator` `protoc` plugin ensures that the given protobuf files contain valid
GAPIC configuration values. The configuration values are supplied via
proto annotations defined at [googleapis/api-common-protos](https://github.com/googleapis/api-common-protos).

Purpose
-------

To provide a tool that emits actionable messages when protobufs are
inproperly configured for GAPIC generation.

Furthermore, to provide a utility that enables GAPIC micro-generator
error message conformance testing.

Installation
------------

To use the plugin, [protoc](https://developers.google.com/protocol-buffers/docs/downloads) must be installed. 

#### Download release binary
```sh
> curl -sSL https://github.com/googleapis/gapic-config-validator/releases/download/v$SEMVER/gapic-config-validator-$SEMVER-$OS-$ARCH.tar.gz | tar xz
> chmod +x protoc-gen-gapic-validator
> export PATH=$PATH:`pwd`
```

#### Via Go tooling
```sh
> go get github.com/googleapis/gapic-config-validator/cmd/protoc-gen-gapic-validator
```

#### From source
```sh
> mkdir -p $GOPATH/src/github.com/googleapis
> cd $GOPATH/src/github.com/googleapis
> git clone https://github.com/googleapis/gapic-config-validator.git
> cd gapic-config-validator
> go install ./cmd/protoc-gen-gapic-validator
```

`make install` executes that last `go install` command for ease of development. 

Invocation
----------

`protoc -I $COMMON_PROTO -I . --gapic-validator_out=. a.proto b.proto`

The `$COMMON_PROTO` variable represents a path to the [googleapis/api-common-protos](https://github.com/googleapis/api-common-protos) directory to import the configuration annotations.

For the time being, the output directory specified by `gapic-validator_out` is not used because there is nothing generated. This value can be anything. 

It is recommended that this validator be invoked prior to `gapic-generator-*` micro-generator invocation.
```sh
protoc -I $COMMON_PROTO \
    -I . \
    --gapic-validator_out=. \
    --go_gapic_out $GO_GAPIC_OUT \
    --go_gapic_opt $GO_GAPIC_OPT
    a.proto b.proto
```

### As a Bazel target

In your WORKSPACE, include the project:
```python
http_archive(
  name = "com_googleapis_gapic_config_validator",
  strip_prefix = "gapic-config-validator-0.2.6",
  urls = ["https://github.com/googleapis/gapic-config-validator/archive/v0.2.6.zip"],
)

load("@com_googleapis_gapic_config_validator//:repositories.bzl", "com_googleapis_gapic_config_validator_repositories")

com_googleapis_gapic_config_validator_repositories()
```

In your BUILD file, configure the target:
```python
load("@com_googleapis_gapic_config_validator//:rules_validate/validate.bzl", "gapic_config_validation")

gapic_config_validation(
  name = "validate_acme_proto",
  srcs = [":acme_proto"],
)
```

The GAPIC v1 config comparison feature can be invoked with the `gapic_yaml` attribute:
```python
gapic_config_validation(
  name = "validate_acme_proto",
  srcs = [":acme_proto"],
  gapic_yaml = ":acme_gapic.yaml"
)
```

_Note: this feature will eventually be removed once the GAPIC v1 config is deprecated fully._

A successful build means there are not issues or discrepancies. A failed build means there
were findings to report, which are found on stderr.

Conformance Testing
-------------------

Micro-generator authors (or other GAPIC config-based plugin authors) can test the conformance of their
error messages against the `gapic-config-validator` using the provided `gapic-error-conformance` testing tool.

The `gapic-error-conformance` utility is a binary that exercises both the `gapic-config-validator` and the targeted
plugin against a set of error mode scenarios. The error emitted by the given plugin is diff'd against
that of the validator and reported to the user. If a plugin error does not conform, the `gapic-error-conformance`
utility will have an exit code of one.

#### Installing `gapic-error-conformance`

##### Download release binary

```sh
> curl -sSL https://github.com/googleapis/gapic-config-validator/releases/download/v$SEMVER/gapic-config-validator-$SEMVER-$OS-$ARCH.tar.gz | tar xz
> chmod +x gapic-error-conformance
> export PATH=$PATH:`pwd`
```

##### Via Go tooling

```sh
> go get github.com/googleapis/gapic-config-validator/cmd/gapic-error-conformance
```

##### From source

```sh
> mkdir -p $GOPATH/src/github.com/googleapis
> cd $GOPATH/src/github.com/googleapis
> git clone https://github.com/googleapis/gapic-config-validator.git
> cd gapic-config-validator
> go install ./cmd/gapic-error-conformance
```

`make install` executes that last `go install` command for ease of development. 

#### Invoking `gapic-error-conformance`

```sh
> gapic-error-conformance -plugin="protoc-gen-go_gapic" -plugin_opts="go-gapic-package=foo.com/bar/v1;bar"
```

##### Options

* `-plugin`: the plugin command to execute. This could the path to an executable or just the
executable itself if it's in the `PATH`.
* `-plugin_opts`: comma-delimited string of options to supply the plugin executable.
* `-verbose`: verbose logging mode. Logs the error messages of the validator and plugin

#### Adding `gapic-error-conformance` scenarios

The scenarios exercised by `gapic-error-conformance` are built into the binary. This means the protobufs
provided as `CodeGeneratorRequest` input are built dynamically. The `scenarios()` method builds
the list of scenarios to exercise. Adding a new scenario means adding the code to build the
protobuf & `CodeGeneratorRequest` here. 

*Note: the proto dependencies required by the GAPIC config annotations are loaded and provided*
*via the `common()` method.*

Testing `gapic-config-validator`
--------------------------------

If you are contributing to this project, run the tests with `make test`.

To view code coverage, run `make coverage`.

#### Test protobuf generation

Some tests require more well-defined descriptors than it makes sense to define by hand in the tests themselves.

The [internal/validator/testdata](/internal/validator/testdata) directory contains protos and their generated types that are used in tests.

Should a change be made to the protos in this directory, the generated types need to be regenerated via `make gen-testdata`. You will need the aforementioned `$COMMON_PROTO` set properly.

Releasing
---------

Follow these steps to make a release:

1. Update the `VERSION` in [release.sh](/release.sh)
2. Open a PR with the **only** `VERSION` bump (notice the prepended `v` for the tag name)
```sh
git add release.sh
git commit -m "release v$VERSION"
```
3. Once version bump PR is merged, create and push the version tag (notice the prepended `v` for the tag name)
```sh
git tag v$VERSION && git push --tags
```
4. Build release assets (Note: must have Docker running for image building)
```sh
make release
```
5. Publish release with `VERSION` tag.   
    a. include the `gapic-config-validator-*.tar.gz ` release assets

    b. push the `latest` and `VERSION` tagged Docker images
  ```sh
  gcloud auth configure-docker
  gcloud docker -- push gcr.io/gapic-images/gapic-config-validator
  gcloud docker -- push gcr.io/gapic-images/gapic-config-validator:$VERSION
  ```
6. (optional) Clean up!
```sh
make clean
```

Disclaimer
----------

This is not an official Google product.