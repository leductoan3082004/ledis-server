GO = go
GO_FMT = gofmt
BINARY_NAME = server

install:
	$(GO) mod tidy

run_server: install
	$(GO) run main.go

format:
	$(GO_FMT) -w .

test:
	$(GO) test -v ./...

build: install
	$(GO) build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)