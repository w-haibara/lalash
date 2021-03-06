lalash: cmd/lalash/main.go *.go */*.go go.mod
	go fmt ./...
	go mod tidy
	go build -o lalash ./cmd/...

.PHONY: run
run:
	go fmt ./...
	go mod tidy
	go run ./cmd/...

.PHONY: init
init:
	go mod init github.com/w-haibara/lalash
	go mod tidy

.PHONY: test
test:
	go fmt ./...
	go test -v

.PHONY: docker
docker:
	docker build -t lalash .
	docker run --rm -it lalash
