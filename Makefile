GO_FILES:=$(shell find . -name '*.go')
NODE_ENV=NODE_ENV=$(if $(DEBUG),development,production)

server: $(GO_FILES)
ifdef DEBUG
	@ # install the non-race version for editor tooling which doesn't parse race
	@ # packages for some reason
	$(GO_ENV) GOBIN=$(TMPDIR) go install ./cmd/... 2>&1 > /dev/null &
endif
	$(GO_ENV) GOBIN=$(GO_BIN_DIR) time go install -v $(if $(DEBUG),-race ,)./cmd/...

server_deps:
	rm -rf vendor
	rm -rf Godeps
	godep save ./...

js_%_watch: 
	$(NODE_ENV) webpack --watch --colors --config js/$*/webpack.js

js_%:
	$(NODE_ENV) webpack --colors --progress --config js/$*/webpack.js
