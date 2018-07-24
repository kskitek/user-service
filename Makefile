# export CGO_ENABLED=0

SVC_NAME = user-service
VERSION = `git rev-parse --short HEAD`

DOCKER_REG = registry.gitlab.com/kskitek/arecar
DOCKER_IMAGE = $(DOCKER_REG)/$(SVC_NAME):$(VERSION)

.PHONY: build build-docker

all: build test run

get-deps:
	go get -t ./...

test:
	go test -race ./...

build:
	go build

run: build
	./$(SVC_NAME)

build-docker: test
	env GOOS=linux go build -o $(SVC_NAME)_linux
	docker build -t $(DOCKER_IMAGE) .

run-docker: build-docker
	docker run -it --rm $(DOCKER_IMAGE)

push: build-docker
	docker push $(DOCKER_IMAGE)

clean:
	rm $(SVC_NAME); rm $(SVC_NAME)_linux
