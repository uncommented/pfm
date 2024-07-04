.PHONY: proto

proto: portfolio/portfolio.proto kis/kis.proto
	protoc --go_opt=paths=source_relative --go_out=. \
		--go-grpc_opt=paths=source_relative --go-grpc_out=. \
		--go-grpc_opt=require_unimplemented_servers=false \
		$^

run: proto
	go run main.go

clean:
	rm -rf **/*.pb.go
