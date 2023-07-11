FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.20 as builder

ARG BUILDPLATFORM="linux/amd64"
ARG GOPATH="/go"
ARG CGO_ENABLED=0
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG Version
ARG GitCommit

WORKDIR ${GOPATH}/src/sample

COPY . .
RUN go mod download
RUN go get -d -v ./...

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -o ${GOPATH}/bin ${GOPATH}/src/sample/helloworld-container/worker
