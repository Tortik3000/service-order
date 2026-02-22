LOCAL_BIN := $(CURDIR)/bin
EASYP_BIN := $(LOCAL_BIN)/easyp
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
GOFUMPT_BIN := $(LOCAL_BIN)/gofumpt
GOIMPORTS_BIN := $(LOCAL_BIN)/goimports
GO_TEST=$(LOCAL_BIN)/gotest
GO_TEST_ARGS=-race -v ./...

INSTALL_CMD = sudo -S  apt update && sudo apt install -y protobuf-compiler

all: generate lint test

lint:
	$(GOFUMPT_BIN) -l -w . || true
	$(GOLANGCI_BIN) run || true

generate: bin-deps .generate build

.PHONY: test
test:
	@echo 'Loading environment variables...'
	@echo $$(cat integration/.env | grep -v '^#' | xargs)

	@echo 'Running tests...'
	export $$(cat integration/.env | grep -v '^#' | xargs) && ${GO_TEST} ${GO_TEST_ARGS}


bin-deps: .bin-deps .install-protoc
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps: .create-bin
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 && \
	go install github.com/rakyll/gotest@v0.0.6 && \
	go install go.uber.org/mock/mockgen@latest && \
	mv $(LOCAL_BIN)/mockgen $(LOCAL_BIN)/mockgen_uber && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1 && \
	go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.11 && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1 && \
	go install golang.org/x/tools/cmd/goimports@v0.19.0 && \
	go install github.com/envoyproxy/protoc-gen-validate@v1.2.1


.install-protoc:
	$(INSTALL_CMD)

.create-bin:
	rm -rf ./bin
	mkdir -p ./bin



.generate:
	$(info Generating code...)

	rm -rf ./generated
	mkdir ./generated

	rm -rf ./docs/spec
	mkdir -p ./docs/spec

	rm -rf ~/.easyp/

	(PATH="$(PATH):$(LOCAL_BIN)" && go generate ./...)
	(PATH="$(PATH):$(LOCAL_BIN)" && $(EASYP_BIN) mod download && $(EASYP_BIN) generate)
	go mod tidy
	$(GOIMPORTS_BIN) -w .


build:
	go mod tidy
	go build -o ./bin/service-order ./cmd/service-order

deploy:
	ansible-playbook -i bot_playbook/inventories/service-order.ini bot_playbook/site.yml