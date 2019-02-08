.PHONY: doc

VERSION=$(shell cat ./VERSION)
BUILDID=$(shell git rev-parse --short HEAD)

all: build

lint-dep:
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@cd $(go env GOPATH)/src/github.com/golangci/golangci-lint/cmd/golangci-lint
	@go install -ldflags "-X 'main.version=$(git describe --tags)' -X 'main.commit=$(git rev-parse --short HEAD)' -X 'main.date=$(date)'"

deps:
	@go mod verify
	@go mod download

doc:
	@godoc -v -http=localhost:6060

build: deps lint
	@go build -v -x -o microservices-keeper .

lint:
	@golangci-lint run

update-version:
	@printf "package main\n\nconst (\n\tVersion=\"$(VERSION)\"\n\tBuildID=\"$(BUILDID)\"\n)\n" > version.go
	@gofmt -s -w version.go
	@git commit -a -m "Version bump"
	@git push

tag:
	@git tag -a v$(VERSION) -m "v$(VERSION)"
	@git push --tags

release: update-version build tag
