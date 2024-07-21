FROM golang:1.22.5-bookworm

RUN apt-get update -y
RUN apt-get install -y unzip
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /tmp
ARG PB_REL="https://github.com/protocolbuffers/protobuf/releases"
ARG PB_VERSION=25.1
ARG arch="x86_64"
RUN curl -LO $PB_REL/download/v$PB_VERSION/protoc-$PB_VERSION-linux-$arch.zip && \
    mkdir -p /proto && \
    unzip protoc-$PB_VERSION-linux-$arch.zip -d /proto && \
    rm protoc-$PB_VERSION-linux-$arch.zip
ENV PATH /proto/bin:$PATH

WORKDIR /pfm
ARG target
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile
COPY proto proto
COPY $target $target

RUN make $target
