.PHONY: proto

proto: portfolio/portfolio.proto
	protoc --go_opt=paths=source_relative --go_out=. \
		--go-grpc_opt=paths=source_relative --go-grpc_out=. \
		--go-grpc_opt=require_unimplemented_servers=false \
		$<

run: proto
	go run server.go

clean:
	rm -rf **/*.pb.go
