GOVERAGE = github.com/haya14busa/goverage
GORELEASER = github.com/goreleaser/goreleaser/v2

os := linux
ifeq ($(shell uname -s), Darwin)
	os := osx
endif

processor := aarch_64
ifeq ($(shell uname -p), x86_64)
	processor := x86_64
endif

PROTOC_VERSION = 29.2
PROTOC = bin/protoc/$(PROTOC_VERSION)/bin/protoc

EXAMPLE_DIR=$(PWD)/examples
DOCS_DIR=$(EXAMPLE_DIR)/doc
PROTOS_DIR=$(EXAMPLE_DIR)/proto

EXAMPLE_CMD=$(PROTOC) --plugin=bin/protoc-gen-doc \
	-Ithirdparty -Itmp/googleapis -Iexamples/proto \
	--doc_out=examples/doc

BOLD = \033[1m
CLEAR = \033[0m
CYAN = \033[36m

help: ## Display this help
	@awk '\
		BEGIN {FS = ":.*##"; printf "Usage: make $(CYAN)<target>$(CLEAR)\n"} \
		/^[a-z0-9]+([\/]%)?([\/](%-)?[a-z\-0-9%]+)*:.*? ##/ { printf "  $(CYAN)%-15s$(CLEAR) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(BOLD)%s$(CLEAR)\n", substr($$0, 5) }' \
		$(MAKEFILE_LIST)

##@: Build

build: ## Build the main binary
	@echo "$(CYAN)Building binary...$(CLEAR)"
	@go build -o bin/protoc-gen-doc ./cmd/protoc-gen-doc

build/examples: $(PROTOC) build tmp/googleapis examples/proto/*.proto examples/templates/*.tmpl ## Build example protos
	@echo "$(CYAN)Making examples...$(CLEAR)"
	@rm -f examples/doc/*
	@$(EXAMPLE_CMD) --doc_opt=docbook,example.docbook:Ignore* examples/proto/*.proto
	@$(EXAMPLE_CMD) --doc_opt=html,example.html:Ignore* examples/proto/*.proto
	@$(EXAMPLE_CMD) --doc_opt=json,example.json:Ignore* examples/proto/*.proto
	@$(EXAMPLE_CMD) --doc_opt=markdown,example.md:Ignore* examples/proto/*.proto
	@$(EXAMPLE_CMD) --doc_opt=examples/templates/asciidoc.tmpl,example.txt:Ignore* examples/proto/*.proto

##@: Test

test/bench: ## Run the bench tests
	@echo "$(CYAN)Running bench tests...$(CLEAR)"
	@go test -bench=.

test/units: fixtures/fileset.pb ## Run unit tests
	@echo "$(CYAN)Running unit tests...$(CLEAR)"
	@go test -cover -race ./ ./cmd/... ./extensions/...

##@: Release

release/snapshot: ## Create a local release snapshot
	@echo "$(CYAN)Creating snapshot build...$(CLEAR)"
	@go run $(GORELEASER) --snapshot --clean

release/validate: ## Run goreleaser checks
	@echo "$(CYAN)Validating release...$(CLEAR)"
	@go run $(GORELEASER) check

##@: Binaries (local installations in ./bin)

$(PROTOC): ## Install protoc compiler
	mkdir -p $(shell dirname $(PROTOC))
	curl -o bin/protoc/$(PROTOC_VERSION)/protoc.zip -LO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(os)-$(processor).zip
	unzip bin/protoc/$(PROTOC_VERSION)/protoc.zip -d bin/protoc/$(PROTOC_VERSION)/

fixtures/fileset.pb: fixtures/*.proto fixtures/generate.go fixtures/nested/*.proto
	@echo "$(CYAN)Generating fixtures...$(CLEAR)"
	@cd fixtures && go generate

tmp/googleapis:
	@echo "$(CYAN)Fetching googleapis...$(CLEAR)"
	@rm -rf tmp/googleapis tmp/protocolbuffers
	@git clone --depth 1 https://github.com/googleapis/googleapis tmp/googleapis
	@rm -rf tmp/googleapis/.git
	@git clone --depth 1 https://github.com/protocolbuffers/protobuf tmp/protocolbuffers
	@cp -r tmp/protocolbuffers/src/* tmp/googleapis/
	@rm -rf tmp/protocolbuffers
