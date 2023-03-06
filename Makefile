VERSION:=$(shell git describe --tags)
LDFLAGS=-s -w
LDFLAGS+=-X "github.com/pot-code/web-cli/cmd.AppVersion=$(VERSION)"

install:
	go install -ldflags '$(LDFLAGS)'
