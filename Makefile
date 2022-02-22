all: format build

build:
	go get ./...
	cd cmd/openrc-exporter/ && go build -v && cd ..

format:
	go fmt ./...

test:
	go test -v

clean:
	go clean ./...

.PHONY: all build clean format test
