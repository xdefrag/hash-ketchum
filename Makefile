PACKAGE := hash-ketchum

PATHDIST := $(CURDIR)/dist/
PATHDISTCLIENT := $(PATHDIST)/client
PATHDISTSERVER := $(PATHDIST)/server

PATHCMD := $(CURDIR)/cmd/
PATHCMDCLIENT := $(PATHCMD)/client
PATHCMDSERVER := $(PATHCMD)/server

PATHBUILD := $(CURDIR)/build
PATHDOCKERCOMPOSE := $(CURDIR)/build/docker-compose.yml

PATHAPIPROTO := ./api/pb/api.proto

PATHADAPTERMOCK := ./pkg/adapter/mock

TESTTAGUNIT := -tags=unit
TESTTAGINTEGRATION := -tags=integration

default: generate test

generate-grpc:
	protoc --go_out=plugins=grpc:. $(PATHAPIPROTO)

generate-mocks:
	go generate $(PATHADAPTERMOCK)

test-unit: generate
	go test $(TESTTAGUNIT) -race ./...

test-integration: generate
	@echo For integration testing redis must be accessable on 0.0.0.0:6379
	go test $(TESTTAGINTEGRATION) -race ./...

build-server: generate
	CGO_ENABLED=0 go build -o $(PATHDISTSERVER) $(PATHCMDSERVER)

build-client: generate
	CGO_ENABLED=0 go build -o $(PATHDISTCLIENT) $(PATHCMDCLIENT)

dc-up:
	docker-compose -f $(PATHDOCKERCOMPOSE) up

redis:
	docker run -p6379:6379 -d redis

generate: generate-grpc generate-mocks

test: test-unit test-integration

build: build-client build-server

docker: docker-server docker-client
