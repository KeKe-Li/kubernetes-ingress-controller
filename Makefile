# goimports
GOFMT_PRIVATE="kubernetes-ingress-controller"
GOFMT_LOCAL="kubernetes-ingress-controller"
GOFMT_FILES=$(shell find . -name '*.go' | grep -v \.pb\.go$ | xargs)

.PHONY: fmt
fmt:
	@goimports -l -w -private "${GOFMT_PRIVATE}" -local "${GOFMT_LOCAL}" ${GOFMT_FILES}

.PHONY: lint
lint:
	@golint ./...

.PHONY: build
build:
	@go build ./...

.PHONY: vet
vet:
	@go vet ./...
