lalash: *.go cmd/lalash/*.go cli/*.go
	gofmt -w *.go cmd/lalash/*.go cli/*.go
	go build ./cmd/...

.PHONY: init
init:
	go mod init lalash
	go mod tidy

.PHONY: test
test:
	gofmt -w *.go cmd/lalash/*.go cli/*.go
	go test -v ./...
