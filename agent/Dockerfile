# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12 
#as builder

RUN apt-get update
RUN apt-get install curl unzip -y

# Install protoc
WORKDIR /usr/local
RUN curl -L -O https://github.com/protocolbuffers/protobuf/releases/download/v3.9.2/protoc-3.9.2-linux-x86_64.zip
RUN unzip protoc-3.9.2-linux-x86_64.zip

# Copy local code to the container image.
WORKDIR /go/src/github.com/googleapis/config-checker
COPY . .

# Install the config validator
RUN go get -u github.com/googleapis/gapic-config-validator/cmd/protoc-gen-gapic-validator

# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go get 
RUN CGO_ENABLED=0 GOOS=linux go build -v -o config-checker

ENV PATH="/usr/local/bin:${PATH}"

ENV GOOGLE_APPLICATION_CREDENTIALS "/go/src/github.com/googleapis/config-checker/credentials.json"

# Run the web service on container startup.
CMD ["/go/src/github.com/googleapis/config-checker/config-checker"]
