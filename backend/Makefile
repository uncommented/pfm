.PHONY: proto

all: portfolio advisor

proto: proto/upbit/upbit.proto proto/kis/kis.proto
	protoc --go_opt=paths=source_relative --go_out=. \
		--go-grpc_opt=paths=source_relative --go-grpc_out=. \
		--go-grpc_opt=require_unimplemented_servers=false \
		$^

portfolio: proto
	go build -o portfolio/main portfolio/*.go

advisor: proto $(wildcard advisor/*.go)
	go build -o advisor/main advisor/*.go

clean:
	rm -rf proto/**/*.pb.go advisor/main portfolio/main
