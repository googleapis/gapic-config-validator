gen-testdata:
	protoc -I ${COMMON_PROTO} -I internal/validator/testdata --go_out=${GOPATH}/src internal/validator/testdata/*.proto

test:
	go test ./...

cover:
	go test ./... -coverprofile=validator.cov
	go tool cover -html=validator.cov

install:
	go install ./cmd/protoc-gen-gapic-validator

clean:
	rm -f validator.cov
	rm -f protoc-gen-gapic-validator

conformance:
	go install ./cmd/gapic-error-conformance
	
image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/protoc-gen-gapic-validator
	docker build -t gcr.io/gapic-images/gapic-config-validator . 
	rm protoc-gen-gapic-validator