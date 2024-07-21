.PHONY: proto

proto: portfolio/upbit/upbit.proto portfolio/kis/kis.proto
	protoc --go_opt=paths=source_relative --go_out=. \
		--go-grpc_opt=paths=source_relative --go-grpc_out=. \
		--go-grpc_opt=require_unimplemented_servers=false \
		$^

portfolio: proto
	go build -o portfolio/pfm_portfolio portfolio/*.go

backend: proto
	go build -o backend/pfm_backend backend/*.go

clean:
	rm -rf portfolio/**/*.pb.go backend/pfm_backend portfolio/pfm_portfolio
