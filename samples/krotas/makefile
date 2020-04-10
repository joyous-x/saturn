# usage:
#	make -B ws=${workspace} target=${target} env=${local dev prod}
# example:
#   make -B ws=. target=krotas env=local

WORKSPACE=$(shell pwd)
ifneq ($(ws), )
WORKSPACE=$(abspath $(ws))
endif

TARGET=$(shell echo $(abspath $(WORKSPACE)) | rev | cut -d "/" -f1 | rev)
ifneq ($(target), ) 
TARGET:=$(target)
endif 

ENV=local
ifneq ($(env), ) 
ENV:=$(env)
endif 

GOPATH=$(abspath $(WORKSPACE)/..)
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


buildDir=$(WORKSPACE)/bin

all: clean config resource build 

build:
	@echo "===> $@ : target=$(TARGET) workspace=$(WORKSPACE)"
	GOFILES=`find $(WORKSPACE) -name '*.go' -not -path '*vendor*'`; GOPATH=$(GOPATH) $(GOFMT) -l -w $$GOFILES;
	if [ "$(ENV)" != "" ]; \
	then \
		CGO_ENABLED=0 GOOS=$(GO_OS) GOPATH=${GOPATH} $(GO) mod tidy; \
		CGO_ENABLED=0 GOOS=$(GO_OS) GOPATH=${GOPATH} $(GO) build -tags netgo -ldflags '$(GO_LDFLAGS)' -o $(buildDir)/$(TARGET); \
	fi
	if [ ! -f "$(buildDir)/$(TARGET)" ]; \
	then \
		echo "===> not exist: $(buildDir)/$(TARGET)" && \
		exit ;\
	fi
	@echo "===> $@ : complete"

config:
	@echo "===> $@ : target=$(TARGET) workspace=$(WORKSPACE)"
	if [ "$(ENV)" != "" ]; then \
		rm -rf $(buildDir)/config && mkdir -p $(buildDir)/config; \
		cp -rf $(WORKSPACE)/env/config/$(ENV)/* $(buildDir)/config; \
	fi
	if [ ! -d "$(buildDir)/config" ]; then \
		echo "===> dir not exist: $(buildDir)/config" && \
		exit ;\
	fi
	cp -rf $(WORKSPACE)/env/crts $(buildDir)/crts;
	cp -rf $(WORKSPACE)/env/nginx.conf $(buildDir)/nginx.conf;
	cp -rf $(WORKSPACE)/env/Dockerfile $(buildDir)/Dockerfile;
	cp -rf $(WORKSPACE)/env/docker-compose.biz.yaml $(buildDir)/docker-compose.biz.yaml;
	cp -rf $(WORKSPACE)/env/docker-compose.env.yaml $(buildDir)/docker-compose.env.yaml;
	if [ "$(ENV)" = "dev" ]; then \
		@echo "this is dev"; \
	fi
	@echo "===> $@ : complete"

rsc_path=$(abspath $(WORKSPACE)/env/rsc)
resource:
	@echo "===> $@ : ready : rsc_path=$(rsc_path)"
	if [ ! -d "$(rsc_path)" ]; then \
		echo "!!! error : $(rsc_path) not exists, please git clone https://github.com/joyous-x/Solar firstly"; \
	else \
		rm -rf $(buildDir)/$@ && mkdir -p $(buildDir)/$@; \
		cp -rf $(rsc_path)/* $(buildDir)/$@; \
	fi
	@echo "===> $@ : complete"

define gofmt_files
	@echo "~~~ start: $1"
	GOFILES=`find $(WORKSPACE) -name '*.go' -not -path '*vendor*'`; \
	GOPATH=$(GOPATH) $(GOFMT) -l -w $$GOFILES;
	@echo "~~~ end: $1"
endef

gofmt:
	$(call gofmt_files,$@) 

.PHONY: all

clean:
	rm -fr $(buildDir)
