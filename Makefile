GO = go
GO_FMT = gofmt
BINARY_NAME = server

run_server:
	$(GO) run main.go

format:
	$(GO_FMT) -w .

test:
	$(GO) test -v ./...

build:
	$(GO) build -o $(BINARY_NAME)