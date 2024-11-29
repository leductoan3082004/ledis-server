GO = go
GO_FMT = gofmt

run_server:
	$(GO) run main.go

format:
	$(GO_FMT) -w .

test:
	$(GO) test -v ./...