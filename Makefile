test:
	go env -w CGO_ENABLED=1
	go test -race ./...