IMAGE=acoshift/cachestatic-backend
TAG=1.1.1
GOLANG_VERSION=1.8
REPO=github.com/acoshift/cachestatic-backend

server: server.go
	go get -v
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o server ./server.go

build:
	docker pull golang:$(GOLANG_VERSION)
	docker run --rm -it -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) golang:$(GOLANG_VERSION) /bin/bash -c "make server"
	docker build --pull -t $(IMAGE):$(TAG) .

push: clean build
	docker push $(IMAGE):$(TAG)

clean:
	rm -f server
