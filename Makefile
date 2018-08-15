# export CGO_ENABLED=0

SVC_NAME = user-service
# VERSION = `git rev-parse --short HEAD`
VERSION = 0.0.1

DOCKER_REG = kskitek
DOCKER_BASE_IMAGE = $(DOCKER_REG)/$(SVC_NAME):$(VERSION)
DOCKER_FLYWAY_IMAGE = $(DOCKER_REG)/$(SVC_NAME)-flyway:$(VERSION)

.PHONY: build build-docker clean

all: verify test build run

get-deps:
	go get -t ./...

verify:
	go vet ./...

test: verify
	go test -race ./...

build:
	go build

run: build
	./$(SVC_NAME)

build-docker: test build
	env GOOS=linux go build -o $(SVC_NAME)_linux
	docker build -t $(DOCKER_BASE_IMAGE) .
	docker build -t $(DOCKER_FLYWAY_IMAGE) -f Dockerfile-flyway .

run-docker:
	docker-compose up

push: build-docker
	docker push $(DOCKER_BASE_IMAGE)
	docker push $(DOCKER_FLYWAY_IMAGE)

clean:
	rm $(SVC_NAME); rm $(SVC_NAME)_linux
