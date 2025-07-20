.PHONY: build test clean run dev-build

dev-build:
	go build -race .

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" .

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" .

build-macos-x86:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" .

build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" .

test:
	go test ./...

run:
	go run .

clean:
	rm -f todo-cli-go todo-cli-go.exe

deps:
	go mod download
	go mod tidy
