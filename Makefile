.PHONY: doc

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
	@go build -v -x -o microservices-keeper main.go

lint:
	@golangci-lint run
