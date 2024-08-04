FROM golang:1.22.5-bookworm

RUN apt-get update -y
RUN apt-get install -y unzip
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
