GOPATH := $(shell go env GOPATH)
GOROOT := $(shell go env GOROOT)
GOBIN := $(GOPATH)/bin

TOOL_GOIMPORTS := $(GOBIN)/goimports
TOOL_GOFMT := $(GOROOT)/bin/gofmt
TOOL_GOLINT := $(GOBIN)/golangci-lint
TOOL_GOVVV := $(GOBIN)/govvv

PKGS := $(shell go list -f '{{.Dir}}' ./...)
LDFLAGS := $(shell $(TOOL_GOVVV) -flags)

BINARY := trollr
VERSION ?= vlatest
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: lint \
		fmt

.DEFAULT_GOAL := run

$(TOOL_GOIMPORTS):
	go get golang.org/x/tools/cmd/goimports
	go mod tidy

$(TOOL_GOLINT):
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
	go mod tidy

$(TOOL_GOVVV):
	go get github.com/JoshuaDoes/govvv
	go mod tidy

lint: $(TOOL_GOLINT)
	$(TOOL_GOLINT) run

fmt: $(TOOL_GOIMPORTS)
	$(TOOL_GOFMT) -w -s $(PKGS)
	$(TOOL_GOIMPORTS) -w $(PKGS)

run:
	go run -ldflags="$(LDFLAGS)" ./app/main.go

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: windows linux darwin