lalash: *.go cmd/lalash/*.go cli/*.go go.mod
	gofmt -w *.go cmd/lalash/*.go cli/*.go
	go mod tidy
	go build ./cmd/...

.PHONY: init
init:
	go mod init lalash
	go mod tidy

.PHONY: test
test:
	gofmt -w *.go cmd/lalash/*.go cli/*.go
	go test -v ./...
