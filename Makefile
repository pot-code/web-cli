VERSION:=$(shell git describe --tags)
GO:=GO111MODULE=on go
LDFLAGS=-s -w
LDFLAGS+=-X "github.com/pot-code/web-cli/cmd.AppVersion=$(VERSION)"

install: generate
	@echo "installing..."
	$(GO) install -ldflags '$(LDFLAGS)'

generate: templates/*.qtpl
	@echo "generating templates"
	$(GO) generate