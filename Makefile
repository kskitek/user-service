# export CGO_ENABLED=0

SVC_NAME = user-service
# VERSION = `git rev-parse --short HEAD`
VERSION = 0.0.1

DOCKER_REG = kskitek
DOCKER_BASE_IMAGE = $(DOCKER_REG)/$(SVC_NAME):$(VERSION)
DOCKER_FLYWAY_IMAGE = $(DOCKER_REG)/$(SVC_NAME)-flyway:$(VERSION)

FLYWAY_VERSION=5.1.4

.PHONY: build build-docker clean

all: verify test build run

get-deps:
	go get -t ./...

verify:
	go vet ./...

test: verify
	go test -race ./...

test-pre-it-docker:
	docker-compose -f docker-compose-it.yaml start
	sleep 5

test-pre-it: clean-test-it test-pre-it-docker flyway-migrate

test-it: test test-pre-it
	env $$(cat config/user-service-it.env) go test -v -race -tags=it ./...

build:
	go build

run:
	env $$(cat config/user-service.env) ./$(SVC_NAME)

build-docker: test build
	env GOOS=linux go build -o $(SVC_NAME)_linux
	docker build -t $(DOCKER_BASE_IMAGE) .
	docker build -t $(DOCKER_FLYWAY_IMAGE) -f Dockerfile-flyway .

run-docker:
	docker-compose up

push: build-docker
	docker push $(DOCKER_BASE_IMAGE)
	docker push $(DOCKER_FLYWAY_IMAGE)

flyway-download:
	wget https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/5.1.4/flyway-commandline-$(FLYWAY_VERSION).tar.gz
	tar -xzvf flyway-commandline-$(FLYWAY_VERSION).tar.gz && rm flyway-commandline-$(FLYWAY_VERSION).tar.gz

flyway-migrate:
	cp -f config/flyway.conf flyway-$(FLYWAY_VERSION)/conf/flyway.conf
	cp -fr sql flyway-$(FLYWAY_VERSION)/sql
	flyway-$(FLYWAY_VERSION)/flyway migrate

clean-test-it:
	docker-compose stop ; echo ""
	docker stop $(SVC_NAME)_pg ; echo ""
	docker container rm $(SVC_NAME)_pg ; echo ""

clean: clean-test-it
	rm $(SVC_NAME); rm $(SVC_NAME)_linux ; echo ""
