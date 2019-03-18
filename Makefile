gen-testdata:
		protoc -I ${COMMON_PROTO} -I internal/validator/testdata --go_out=${GOPATH}/src internal/validator/testdata/*.proto

test:
	go test ./...

cover:
	go test ./... -coverprofile=validator.cov
	go tool cover -html=validator.cov

install:
	go install ./cmd/protoc-gen-gapic-validator