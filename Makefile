GOPATH := $(shell go env GOPATH)
GOROOT := $(shell go env GOROOT)
GOBIN := $(GOPATH)/bin

TOOL_GOIMPORTS := $(GOBIN)/goimports
TOOL_GOFMT := $(GOROOT)/bin/gofmt
TOOL_GOLINT := $(GOBIN)/golangci-lint
TOOL_GOVVV := $(GOBIN)/govvv
TOOL_SG2MD := $(shell npm bin -g)/swagger-markdown

PKGS := $(shell go list -f '{{.Dir}}' ./...)

BINARY := trollr
VERSION ?= vlatest
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: \
		$(PLATFORMS) \
		build \
		clean \
		deps \
		fmt \
		lint \
		run \
		test \
        container \
        release \
        run-container \

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

$(TOOL_SG2MD):
	npm install -g swagger-markdown

clean:
	rm -rf out

lint: $(TOOL_GOLINT)
	$(TOOL_GOLINT) run

fmt: $(TOOL_GOIMPORTS)
	$(TOOL_GOFMT) -w -s $(PKGS)
	$(TOOL_GOIMPORTS) -w $(PKGS)

deps: $(TOOL_GOVVV) $(TOOL_GOLINT) $(TOOL_GOIMPORTS)
	go get -v -t -d ./...

test:
	go test -v ./...

run: $(TOOL_GOVVV)
	go run -v -ldflags="$(shell $(TOOL_GOVVV) -flags)" ./app

build: $(TOOL_GOVVV)
	mkdir -p out/build
	go build -v -ldflags="$(shell $(TOOL_GOVVV) -flags)" -o out/build/trollr ./app

swagger:
	docker run --rm -it -e GOPATH=$$HOME/go:/go -v $$HOME:$$HOME -w $$(pwd) quay.io/goswagger/swagger generate spec -o static/swagger.json --scan-models
	$(TOOL_SG2MD) -i ./static/swagger.json -o API.md

$(PLATFORMS): $(TOOL_GOVVV)
	mkdir -p out/release
	GOOS=$(os) GOARCH=amd64 CGO_ENABLED=0 go build -v -tags netgo -ldflags "-extldflags \"-static\" -s -w $(shell $(TOOL_GOVVV) -flags)" -a -o out/release/$(BINARY)-$(VERSION)-$(os)-amd64 ./app

release: windows linux darwin

container:
	docker build --tag bendoerr/trollr:latest .

run-container: container
	docker run --publish 7891:7891 --interactive --tty --label trollr bendoerr/trollr:latest