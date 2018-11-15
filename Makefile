PACKAGE := hash-ketchum

PATHDIST := $(CURDIR)/dist/
PATHDISTCLIENT := $(PATHDIST)/client
PATHDISTSERVER := $(PATHDIST)/server

PATHCMD := $(CURDIR)/cmd/
PATHCMDCLIENT := $(PATHCMD)/client
PATHCMDSERVER := $(PATHCMD)/server

PATHBUILD := $(CURDIR)/build
PATHDOCKERCOMPOSE := $(CURDIR)/build/docker-compose.yml

PATHAPIPROTO := $(CURDIR)/api/pb/api.proto

PATHADAPTERMOCK := $(CURDIR)/pkg/adapter/mock

TESTTAGUNIT := -tags=unit
TESTTAGINTEGRATION := -tags=integration

default: build

generate-grpc:
	protoc --go_out=plugins=grpc:. $(PATHAPIPROTO)

generate-mocks:
	go generate $(PATHADAPTERMOCK)

test-unit:
	go test $(TESTTAGUNIT) -race ./...

test-integration:
	@echo For integration testing redis must be accessable on 0.0.0.0:6379
	go test $(TESTTAGINTEGRATION) -race ./...

build-server:
	CGO_ENABLED=0 go build -o $(PATHDISTSERVER) $(PATHCMDSERVER)

build-client:
	CGO_ENABLED=0 go build -o $(PATHDISTCLIENT) $(PATHCMDCLIENT)

dc-up:
	docker-compose -f $(PATHDOCKERCOMPOSE) up

generate: generate-grpc generate-mocks

test: test-unit test-integration

build: build-client build-server

docker: docker-server docker-client
