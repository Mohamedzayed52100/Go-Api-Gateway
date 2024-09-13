SHELL := /bin/bash

.PHONY: gen*

# detect if the CI mode is enabled
IS_CI_ENABLED := $(shell [ -z "$$CI" ] && v=false || v=true; echo $$v)

OS := $(shell uname)
OS_ALT := $(shell if [ "$(OS)" == "Darwin" ]; then echo "osx"; else echo $(OS); fi )
ARCH := $(shell uname -m)
ARCH_ALT := $(shell if [ "$(ARCH)" == "arm64" ]; then echo "aarch_64"; else echo $(ARCH); fi )
IMAGE_NAME := "goplace/openapi-generator:latest"
PROTOC_IMAGE := "namely/protoc-all:latest"
GOPATH := $(shell go env GOPATH)
PROTO_REPOS := "github.com/goplaceapp/goplace-guest"


## START API/Proto generator targets ##
gen:
	$(MAKE) gen-openapi
	#$(MAKE) gen-proto

gen-openapi:
	docker run --rm -v $$(pwd):/local $(IMAGE_NAME) -- -r /local -f ./openapi/specs/openapi.yaml -o ./openapi

#gen-proto:
#    @for r in $(PROTO_REPOS); do \
#        git clone https://$$r.git /tmp/proto-repo && \
#        docker run --rm \
#            -v /tmp/proto-repo:/defs \
#            -v $(GOPATH)/src:/go/src \
#            $(PROTOC_IMAGE) \
#            --go_out=/go/src \
#            --go_opt=paths=source_relative \
#            --go-grpc_out=require_unimplemented_servers=false:/go/src \
#            --go-grpc_opt=paths=source_relative; \
#        rm -rf /tmp/proto-repo; \
#	done

clean:
	find ./openapi -type f -name '*.go' | xargs rm -f
## END API/Proto generator targets ##
