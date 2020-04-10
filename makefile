# usage:
#	make -B ws=${ws}
# example:
#   make -B ws=.

WORKSPACE=$(shell pwd)
ifneq ($(ws), )
WORKSPACE=$(abspath $(ws))
endif

GOPATH=$(abspath $(WORKSPACE))
GO=go
GO_OS ?= $(shell go env GOOS)
GOFMT=gofmt
GODEP=dep

BUILD_TIME=$(shell date '+%Y-%m-%dT%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null)
GO_LDFLAGS:= -X internal/common/version.BuildTime=$(BUILD_TIME)
GO_LDFLAGS+= -X internal/common/version.Ver=${CI_COMMIT_TAG}
GO_LDFLAGS+= -X internal/common/version.GitVersion=$(GIT_COMMIT)
GO_LDFLAGS+= -extldflags -static

gofmt:
	GOFILES=`find $(WORKSPACE) -name '*.go' -not -path '*vendor*' -not -path '*pkg/mod/*'`; GOPATH=$(GOPATH) $(GOFMT) -l -w $$GOFILES;