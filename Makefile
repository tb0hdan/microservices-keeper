.PHONY: doc

all: build

doc:
	@godoc -v -http=localhost:6060

build:
	@go build -v -x -o microservices-keeper main.go

