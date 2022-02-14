GOPATH := $(shell go env GOPATH)

.PHONY: dev
dev:
	$(GOPATH)/bin/air -c .air.toml