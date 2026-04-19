BINARY     := portwatch
CMD        := ./cmd/portwatch
GO         := go
GOFLAGS   ?=

.PHONY: all build test lint clean install

all: build

build:
	$(GO) build $(GOFLAGS) -o bin/$(BINARY) $(CMD)

test:
	$(GO) test ./...

test-short:
	$(GO) test -short ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found"; exit 1; }
	golangci-lint run ./...

install:
	$(GO) install $(GOFLAGS) $(CMD)

clean:
	rm -rf bin/
