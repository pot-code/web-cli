VERSION:=$(shell git describe --tags)
GO:=GO111MODULE=on go
LDFLAGS=-s -w
LDFLAGS+=-X "github.com/pot-code/web-cli/cmd.AppVersion=$(VERSION)"

install: 
	$(GO) install -ldflags '$(LDFLAGS)'
