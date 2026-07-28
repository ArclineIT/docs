BINARY := docs-server
GOFLAGS := -trimpath -ldflags="-s -w"

.PHONY: build run linux-amd64 linux-arm64 all test clean

build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o $(BINARY) ./cmd/serve

run:
	go run ./cmd/serve

linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BINARY)-linux-amd64 ./cmd/serve

linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o $(BINARY)-linux-arm64 ./cmd/serve

all: linux-amd64 linux-arm64

test:
	go test ./...

clean:
	rm -f $(BINARY) $(BINARY)-linux-*

# Static site build
.PHONY: build-static
build-static:
	go run ./cmd/build

build-static-watch:
	go run ./cmd/build --watch


