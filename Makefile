list := *.go

lalash: $(list) go.mod
	gofmt -w $(list)
	go mod tidy
	go build ./...

.PHONY: init
init:
	go mod init lalash
	go mod tidy

.PHONY: test
test:
	gofmt -w $(list)
	go test -v ./...

.PHONY: docker
docker:
	docker build -t lalash .
	docker run --rm -it lalash
