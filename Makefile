test:
	go env -w CGO_ENABLED=1
	go test -race ./...

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/urlshortener.proto