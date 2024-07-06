.PHONY: proto

proto: upbit/upbit.proto kis/kis.proto
	protoc --go_opt=paths=source_relative --go_out=. \
		--go-grpc_opt=paths=source_relative --go-grpc_out=. \
		--go-grpc_opt=require_unimplemented_servers=false \
		$^

server: proto
	go build -o server/server server/server.go

client: proto
	go build -o client/client client/client.go

clean:
	rm -rf **/*.pb.go client/client server/server
