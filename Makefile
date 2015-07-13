DEBUG_FLAG = $(if $(DEBUG),-v)
GOPATH_ENV="$(PWD)/.godeps:$(PWD)"
GOBIN_ENV="$(PWD)/.godeps/bin"
GOPKG_ENV="$(PWD)/.godeps/pkg/darwin_amd64/github.com/yoheimuta/dbq"

deps:
	wget -qO- https://raw.githubusercontent.com/pote/gpm/v1.2.3/bin/gpm | GOPATH=$(GOPATH_ENV) bash

check:
	goimports -l main.go command/
	vet main.go; vet command/
	golint ./...

fix:
	goimports -w main.go command/

test:
	GOPATH=$(GOPATH_ENV) go test $(DEBUG_FLAG) ./...

install:
	rm -r $(GOPKG_ENV); GOBIN=$(GOBIN_ENV) GOPATH=$(GOPATH_ENV) go install
