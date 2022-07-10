
### NAMES AND LOCS ############################
APPNAME      = shinpuru
PACKAGE      = github.com/zekroTJA/shinpuru
LDPAKAGE     = internal/util
CONFIG       = $(CURDIR)/config/private.config.yml
BINPATH      = $(CURDIR)/bin
PRETTIER_CFG = "$(CURDIR)/.prettierrc.yml"
TMPBIN       = "./bin/tmp/$(APPNAME)"
###############################################

### EXECUTABLES ###############################
GO     	      = go
GOLINT 	      = golint
GREP   	      = grep
NPM    	      = npm
PRETTIER      = prettier
YARN					= yarn
DOCKERCOMPOSE = docker-compose
SWAGGO		    = swag
SWAGGER2MD    = swagger-markdown
###############################################

# ---------------------------------------------

BIN = $(BINPATH)/$(APPNAME)

TAG        = $(shell git describe --tags)
COMMIT     = $(shell git rev-parse HEAD)
DATE       = $(shell date +%s)


ifneq ($(GOOS),)
	BIN := $(BIN)_$(GOOS)
endif

ifneq ($(GOARCH),)
	BIN := $(BIN)_$(GOARCH)
endif

ifneq ($(TAG),)
	BIN := $(BIN)_$(TAG)
endif

ifeq ($(OS),Windows_NT)
	ifeq ($(GOOS),)
		BIN := $(BIN).exe
	endif
endif

ifeq ($(GOOS),windows)
	BIN := $(BIN).exe
endif


PHONY = _make
_make: deps fe copyfe build cleanup

PHONY += build
build: $(BIN)

PHONY += deps
deps:
	$(GO) mod tidy
	cd ./web && \
		$(NPM) install

$(BIN):
	$(GO) build  \
		-v -o $@ -ldflags "\
			-X $(PACKAGE)/$(LDPAKAGE).AppVersion=$(TAG) \
			-X $(PACKAGE)/$(LDPAKAGE).AppCommit=$(COMMIT) \
			-X $(PACKAGE)/$(LDPAKAGE).AppDate=$(DATE) \
			-X $(PACKAGE)/$(LDPAKAGE).Release=TRUE" \
		$(CURDIR)/cmd/$(APPNAME)/*.go

PHONY += test
test:
	$(GO) test -race -v -cover ./...

PHONY += lint
lint:
	$(GOLINT) ./... | $(GREP) -v vendor || true

$(TMPBIN):
	$(GO) build -race -v -o $@ $(CURDIR)/cmd/$(APPNAME)/*.go

PHONY += run
run: $(TMPBIN)
	$(TMPBIN) -c $(CONFIG) -quiet

PHONY += rundev
rundev: $(TMPBIN)
	$(TMPBIN) -devmode -c $(CONFIG) -quiet

PHONY += cleanup
cleanup:

PHONY += fe
fe:
	cd $(CURDIR)/web && \
		$(YARN) run build

PHONY += copyfe
copyfe:
	cp -R web/dist/web/* internal/util/embedded/webdist

PHONY += runfe
runfe:
	cd ./web && $(YARN) start

PHONY += prettify
prettify:
	$(PRETTIER) \
	    --config $(PRETTIER_CFG) \
	    --write \
	    	$(CURDIR)/web/src/**/*.js \
	    	$(CURDIR)/web/src/**/**/*.js \
	    	$(CURDIR)/web/src/**/*.vue \
	    	$(CURDIR)/web/src/**/**/*.vue

PHONY += devstack
devstack:
	$(DOCKERCOMPOSE) -f docker-compose.dev.yml \
		up -d

PHONY += devstack
APIDOCS_OUTDIR = "$(CURDIR)/docs/restapi/v1"
apidocs:
	$(SWAGGO) init \
		-g $(CURDIR)/internal/services/webserver/v1/router.go \
		-o $(APIDOCS_OUTDIR) \
		--parseDependency --parseDepth 2
	rm $(APIDOCS_OUTDIR)/docs.go
	$(SWAGGER2MD) -i $(APIDOCS_OUTDIR)/swagger.json -o $(APIDOCS_OUTDIR)/restapi.md

PHONY += help
help:
	@echo "Available targets:"
	@echo "  #        - creates binary in ./bin"
	@echo "  cleanup  - tidy up temporary stuff created by build or scripts"
	@echo "  deps     - ensure dependencies are installed"
	@echo "  devstack - spins up the dev docker-compose stack"
	@echo "  fe       - build font end files"
	@echo "  lint     - run linters (golint)"
	@echo "  run      - debug run app (go run) with test config"
	@echo "  runfe    - debug run front end vue live-server"
	@echo "  test     - run tests (go test)"
	@echo ""
	@echo "Cross Compiling:"
	@echo "  (env GOOS=linux GOARCH=arm make)"
	@echo ""
	@echo "Use different configs for run:"
	@echo "  make CONF=./myCustomConfig.yml run"
	@echo ""


.PHONY: $(PHONY)
