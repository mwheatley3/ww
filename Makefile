GO_FILES:=$(shell find . -name '*.go')
	
server: $(GO_FILES)
ifdef DEBUG
	@ # install the non-race version for editor tooling which doesn't parse race
	@ # packages for some reason
	$(GO_ENV) GOBIN=$(TMPDIR) go install ./cmd/... 2>&1 > /dev/null &
endif
	$(GO_ENV) GOBIN=$(GO_BIN_DIR) time go install -v $(if $(DEBUG),-race ,)./cmd/...
