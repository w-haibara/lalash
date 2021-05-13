list := *.go

lalash: $(list) go.mod
	gofmt -w $(list)
	go mod tidy
	go build -o lalash 

.PHONY: init
init:
	go mod init lalash
	go mod tidy

.PHONY: docker
docker:
	docker build -t lalash .
	docker run --rm -it lalash
