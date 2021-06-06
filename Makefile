lalash: cmd/lalash/main.go *.go */*.go go.mod
	go fmt ./...
	go mod tidy
	go build -o lalash ./cmd/...

.PHONY: run
run:
	go fmt ./...
	go mod tidy
	go run .

.PHONY: init
init:
	go mod init lalash
	go mod tidy

.PHONY: docker
docker:
	docker build -t lalash .
	docker run --rm -it lalash
