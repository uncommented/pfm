.PHONY: proto

all: portfolio advisor

proto: proto/upbit/upbit.proto proto/kis/kis.proto
	protoc --go_opt=paths=source_relative --go_out=. \
		--go-grpc_opt=paths=source_relative --go-grpc_out=. \
		--go-grpc_opt=require_unimplemented_servers=false \
		$^

portfolio: proto
	go build -o portfolio/pfm_portfolio portfolio/*.go

advisor: proto $(wildcard advisor/*.go)
	go build -o pfm_backend advisor/*.go

clean:
	rm -rf proto/**/*.pb.go advisor/pfm_backend portfolio/pfm_portfolio
