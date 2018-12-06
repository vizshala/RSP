default: build

build: all

all: 
	go install ./service/...

test: test-all

test-all:
	@go test -v ./...