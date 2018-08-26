# export CGO_ENABLED=0

SVC_NAME = user-service
# VERSION = `git rev-parse --short HEAD`
VERSION = 0.10.0

DOCKER_REG = kskitek
DOCKER_BASE_IMAGE = $(DOCKER_REG)/$(SVC_NAME):$(VERSION)
DOCKER_FLYWAY_IMAGE = $(DOCKER_REG)/$(SVC_NAME)-flyway:$(VERSION)

FLYWAY_VERSION=5.1.4

.PHONY: build build-docker clean

all: help


get-flyway:
	wget https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/5.1.4/flyway-commandline-$(FLYWAY_VERSION).tar.gz
	tar -xzvf flyway-commandline-$(FLYWAY_VERSION).tar.gz && rm flyway-commandline-$(FLYWAY_VERSION).tar.gz

## get-deps: fetches dependencies of project
get-deps: get-flyway
	go get -t ./...

verify: build
	go vet ./...

## test: tests go service
test: verify
	go test -race ./...

test-pre-it-docker:
	docker-compose -f docker-compose-it.yaml start
	sleep 5

## test-pre-it: sets up integration tests environment (based on docker-compose-it.yaml)
test-pre-it: clean-test-it test-pre-it-docker migrate

## test-it: integration tests
test-it: test test-pre-it
	env $$(cat config/user-service-it.env) go test -v -race -tags=it ./...

## build: build go service
build:
	go build

## run: runs service locally
run:
	env $$(cat config/user-service.env) ./$(SVC_NAME)

## build-docker: builds docker image of the service
build-docker: test build
	env GOOS=linux go build -o $(SVC_NAME)_linux
	docker build -t $(DOCKER_BASE_IMAGE) .
	docker build -t $(DOCKER_FLYWAY_IMAGE) -f Dockerfile-flyway .

## run-docker: runs service in docker together with required infrastrucure containers
run-docker:
	docker-compose up

## push: pushes service image to docker registry
push: build-docker
	docker push $(DOCKER_BASE_IMAGE)
	docker push $(DOCKER_FLYWAY_IMAGE)

## migrate: migrates database specified in config/flyway.conf
migrate:
	cp -f config/flyway.conf flyway-$(FLYWAY_VERSION)/conf/flyway.conf
	cp -fr sql flyway-$(FLYWAY_VERSION)/sql
	flyway-$(FLYWAY_VERSION)/flyway migrate

## clean-test-it: cleans integration tests environment
clean-test-it:
	docker-compose stop ; echo ""
	docker stop $(SVC_NAME)_pg ; echo ""
	docker container rm $(SVC_NAME)_pg ; echo ""

## clean: cleans service binary files and integration tests environment
clean: clean-test-it
	rm $(SVC_NAME); rm $(SVC_NAME)_linux ; echo ""

help: Makefile
	@echo " Choose a command run in \033[32m"$(SVC_NAME)"\033[0m:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
