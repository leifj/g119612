MODULE = $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
PACKAGES := $(shell go list ./... | grep -v /vendor/)
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

.PHONY: default
default: build

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test:
	go test ./pkg/etsi119612

.PHONY: build
build:  ## build the library
	CGO_ENABLED=0 go build ${LDFLAGS} -o etsi_ts -a cmd/etsi_ts/main.go

.PHONY: clean
clean: ## remove temporary files
	go clean

.PHONY: realclean
realclean: ## remove generated files - requires "make gen"
	rm -f pkg/etsi119612/*.xsd.go

.PHONY: gen
gen: ## generate code from xsd
	xgen -i xsd -o pkg/etsi119612 -l Go -p etsi119612
