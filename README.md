gapic-config-validator
======================

A configuration validator for GAPIC config proto annotations.

This `protoc` plugin ensures that the given protobuf files contain valid
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

Conformance Testing
-------------------

Micro-generator authors (or other GAPIC config-based plugin authors) can test the conformance of their
error messages against the `gapic-config-validator` using the provided `conformance` testing tool.

The `conformance` utility is a binary that exercises both the `gapic-config-validator` and the targeted
plugin against a set of error mode scenarios. The error emitted by the given plugin is diff'd against
that of the validator and reported to the user. If a plugin error does not conform, the `conformance`
utility will have an exit code of one.

#### Installing `conformance`

##### Download release binary

```sh
> curl -sSL https://github.com/googleapis/gapic-config-validator/releases/download/v$SEMVER/gapic-config-validator-$SEMVER-$OS-$ARCH.tar.gz | tar xz
> chmod +x conformance
> export PATH=$PATH:`pwd`
```

##### Via Go tooling

```sh
> go get github.com/googleapis/gapic-config-validator/cmd/conformance
```

##### From source

```sh
> mkdir -p $GOPATH/src/github.com/googleapis
> cd $GOPATH/src/github.com/googleapis
> git clone https://github.com/googleapis/gapic-config-validator.git
> cd gapic-config-validator
> go install ./cmd/conformance
```

`make install` executes that last `go install` command for ease of development. 

#### Invoking `conformance`

```sh
> conformance -plugin="protoc-gen-go_gapic" -plugin_opts="go-gapic-package=foo.com/bar/v1;bar"
```

##### Options

* `-plugin`: the plugin command to execute. This could the path to an executable or just the
executable itself if it's in the `PATH`.
* `-plugin_opts`: comma-delimited string of options to supply the plugin executable.
* `-verbose`: verbose logging mode. Logs the error messages of the validator and plugin

#### Adding `conformance` scenarios

The scenarios exercised by `conformance` are built into the binary. This means the protobufs
provided as `CodeGeneratorRequest` input are built dynamically. The `scenarios()` method builds
the list of scenarios to exercise. Adding a new scenario means adding the code to build the
protobuf & `CodeGeneratorRequest` here. 

*Note: the proto dependencies required by the GAPIC config annotations are loaded and provided*
*via the `common()` method.*

Testing `gapic-config-validator`
--------------------------------

If you are contributing to this project, run the tests with `make test`.

#### Test protobuf generation

Some tests require more well-defined descriptors than it makes sense to define by hand in the tests themselves.

The [internal/validator/testdata](/internal/validator/testdata) directory contains protos and their generated types that are used in tests.

Should a change be made to the protos in this directory, the generated types need to be regenerated via `make gen-testdata`. You will need the aforementioned `$COMMON_PROTO` set properly.

Disclaimer
----------

This is not an official Google product.