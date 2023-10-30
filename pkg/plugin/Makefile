.PHONY: all
all: build-plugins

.PHONY: build-plugins
build-plugins:
	CGO_ENABLED=$(CGO_ENABLED) sh -c "'$(CURDIR)/scripts/plugin.sh'"

.PHONY: clean
clean:
	-rm ./bin/manager

.PHONY: fmt
fmt: import ## Run go fmt against code.
	go fmt ./...

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

CHECK_LINT?= $(LOCALBIN)/golangci-lint
CHECK_IMPORTS?= $(LOCALBIN)/goimports

.PHONY: lint
lint: check-lint
	GOBIN=$(LOCALBIN) CGO_ENABLED=$(CGO_ENABLED) golangci-lint run -v --timeout=5m ./...

.PHONY: check-lint
check-lint: $(CHECK_LINT) ## Download golangci-lint-setup locally if necessary.
$(CHECK_LINT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) CGO_ENABLED=$(CGO_ENABLED) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: import
import: check-goimports
	GOBIN=$(LOCALBIN) CGO_ENABLED=$(CGO_ENABLED) ./bin/goimports -w -e -l .

.PHONY: check-goimports
check-goimports: $(CHECK_IMPORTS) ## Download goimports-setup locally if necessary.
$(CHECK_IMPORTS): $(LOCALBIN)
	GOBIN=$(LOCALBIN) CGO_ENABLED=$(CGO_ENABLED) go install golang.org/x/tools/cmd/goimports@latest
