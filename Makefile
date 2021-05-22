list := main.go

lalash: $(list) go.mod
	go fmt ./...
	go mod tidy
	go build -o lalash 

.PHONY: run
run:
	go fmt ./...
	go mod tidy
	go run $(list)

.PHONY: init
init:
	go mod init lalash
	go mod tidy

.PHONY: docker
docker:
	docker build -t lalash .
	docker run --rm -it lalash
