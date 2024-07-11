.PHONY: proto

proto: upbit/upbit.proto kis/kis.proto
	protoc --go_opt=paths=source_relative --go_out=. \
		--go-grpc_opt=paths=source_relative --go-grpc_out=. \
		--go-grpc_opt=require_unimplemented_servers=false \
		$^

server: proto
	go build -o server/server server/*.go

backend: proto
	go build -o backend/pfm_backend backend/*.go

clean:
	rm -rf **/*.pb.go backend/pfm_backend server/server
