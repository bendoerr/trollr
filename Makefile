GOPATH := $(shell go env GOPATH)
GOROOT := $(shell go env GOROOT)
GOBIN := $(GOPATH)/bin

TOOL_GOIMPORTS := $(GOBIN)/goimports
TOOL_GOFMT := $(GOROOT)/bin/gofmt
TOOL_GOLINT := $(GOBIN)/golangci-lint
TOOL_GOVVV := $(GOBIN)/govvv

PKGS := $(shell go list -f '{{.Dir}}' ./...)

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

run: $(TOOL_GOVVV)
	go run -ldflags="$(shell $(TOOL_GOVVV) -flags)" ./app/main.go

.PHONY: $(PLATFORMS)
$(PLATFORMS): $(TOOL_GOVVV)
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -ldflags "-extldflags \"-static\" -s -w $(shell $(TOOL_GOVVV) -flags)" -a -o release/$(BINARY)-$(VERSION)-$(os)-amd64 ./app/main.go

.PHONY: release
release: windows linux darwin

.PHONY: container
container:
	docker build --tag bendoerr/trollr:latest .

.PHONY: run-container
run-container: container
	docker run --publish 7891:7891 --interactive --tty --label trollr bendoerr/trollr:latest