PACKAGE  = <%=pkgName%>
DATE     = $(shell date +%s)
BIN      = $(GOPATH)/bin
BASE     = $(GOPATH)/src/<%=repoUrl%>
PKGS     = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))

COMPILED      ?= $(shell date +%FT%T%z)
PATCH_VERSION ?= $(shell if [ -z "$(BUILD_TAG)" ] && [ "$(BUILD_TAG)xxx" == "xxx" ]; then echo "dev"; else echo $(BUILD_TAG); fi)
GIT_COMMIT    ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo unknown)

GO       = go
GODOC    = godoc
GOFMT    = gofmt
DOCKER   = docker
TIMEOUT  = 15

LD_FLAGS = "-X <%=repoUrl%>.PatchVersion=$(PATCH_VERSION) -X <%=repoUrl%>.Compiled=$(DATE) -X <%=repoUrl%>.GitCommit=$(GIT_COMMIT)"
DOCKER_BUILD_ARGS = '--build-arg LD_FLAGS=${LD_FLAGS}'

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt gen lint build

# Tools

GOLINT = $(BIN)/golint
$(BIN)/golint: | $(BASE) ; $(info $(M) installing golint…)
	$Q $(GO) get github.com/golang/lint/golint

GOCOVMERGE = $(BIN)/gocovmerge
$(BIN)/gocovmerge: | $(BASE) ; $(info $(M) installing gocovmerge…)
	$Q $(GO) get github.com/wadey/gocovmerge

GOCOV = $(BIN)/gocov
$(BIN)/gocov: | $(BASE) ; $(info $(M) installing gocov…)
	$Q $(GO) get github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: | $(BASE) ; $(info $(M) installing gocov-xml…)
	$Q $(GO) get github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: | $(BASE) ; $(info $(M) installing go2xunit…)
	$Q $(GO) get github.com/tebeka/go2xunit

GODEP = $(BIN)/dep
$(BIN)/dep: | $(BASE) ; $(info $(M) installing dep…)
	$Q $(GO) get -u github.com/golang/dep/cmd/dep

PROTOC_GRPC_GATEWAY = $(BIN)/protoc-gen-grpc-gateway
$(BIN)/protoc-gen-grpc-gateway: | $(BASE) ; $(info $(M) installing protoc-gen-grpc-gateway…)
	$Q $(GO) get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

PROTOC_SWAGGER = $(BIN)/protoc-gen-swagger
$(BIN)/protoc-gen-swagger: | $(BASE) ; $(info $(M) installing protoc-gen-swagger…)
	$Q $(GO) get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

PROTOC_GO = $(BIN)/protoc-gen-go
$(BIN)/protoc-gen-go: | $(BASE) ; $(info $(M) installing protoc-gen-go…)
	$Q $(GO) get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: install-dependencies 
install-dependencies: $(BIN)/golint $(BIN)/gocovmerge $(BIN)/gocov $(BIN)/gocov-xml $(BIN)/go2xunit $(BIN)/dep $(BIN)/protoc-gen-grpc-gateway $(BIN)/protoc-gen-swagger $(BIN)/protoc-gen-go   ## Install dependent go tools

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: fmt lint vendor | $(BASE) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q cd $(BASE) && $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint vendor | $(BASE) $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q cd $(BASE) && 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint vendor test-coverage-tools | $(BASE) ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)/coverage
	$Q cd $(BASE) && for pkg in $(TESTPKGS); do \
		$(GO) test \
			-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
					grep '^$(PACKAGE)/' | grep -v '^$(PACKAGE)/vendor/' | \
					tr '\n' ',')$$pkg \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	 done
	$Q $(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: lint
lint: vendor | $(BASE) $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

# Dependency management

.PHONY: vendor
vendor: 
	$Q cd $(BASE) && $(GODEP) ensure -v; $(info $(M) retrieving dependencies…)

.PHONY: gen
gen: 
	$Q cd $(BASE) && $(GO) generate ./pkg/api ./pkg/state; $(info $(M) generating api and metrics & trace store …)

.PHONY: build
build: gen  ## Build service
	$Q cd $(BASE) && $(GO) build -ldflags $(LD_FLAGS) \
		-o bin/$(PACKAGE) ./cmd/<%=serviceName%> ; $(info $(M) building executable…)

.PHONY: docker-build
docker-build: ## Build docker image
	$Q cd $(BASE) && $(DOCKER) build '$(DOCKER_BUILD_ARGS)' -t <%=pkgName%>:$(PATCH_VERSION) . && \
		$(DOCKER) tag <%=pkgName%>:$(PATCH_VERSION) <%=containerRegistry%>/<%=repoName%>/<%=pkgName%>:$(PATCH_VERSION); $(info $(M) building docker image…)

.PHONY: docker-push
docker-push: docker-build ## Build and publish docker image
	$Q $(DOCKER) push <%=containerRegistry%>/<%=repoName%>/<%=pkgName%>:$(PATCH_VERSION); $(info $(M) pushing docker image…)

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'