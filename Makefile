DEBUG_FLAG = $(if $(DEBUG),-v)
GOPATH_ENV="$(PWD)/.godeps:$(PWD)"
GOBIN_ENV="$(PWD)/.godeps/bin"

deps:
	wget -qO- https://raw.githubusercontent.com/pote/gpm/v1.2.3/bin/gpm | GOPATH=$(GOPATH_ENV) bash

check:
	goimports -l main.go src/
	vet main.go; vet src/
	golint ./...

fix:
	goimports -w main.go src/

test:
	GOPATH=$(GOPATH_ENV) go test $(DEBUG_FLAG) ./...

install:
	GOBIN=$(GOBIN_ENV) GOPATH=$(GOPATH_ENV) go install
