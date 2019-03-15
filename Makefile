gen-testdata:
		protoc -I ${COMMON_PROTO} -I internal/validator/testdata --go_out=${GOPATH}/src internal/validator/testdata/*.proto

test:
	go test ./...

install:
	go install ./cmd/protoc-gen-gapic-validator