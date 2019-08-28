    
### NAMES AND LOCS ############################
APPNAME      = shinpuru
PACKAGE      = github.com/zekroTJA/shinpuru
LDPAKAGE     = internal/util
CONFIG       = $(CURDIR)/config/private.config.yml
BINPATH      = $(CURDIR)/bin
PRETTIER_CFG = "$(CURDIR)/.prettierrc.yml"
###############################################

### EXECUTABLES ###############################
GO     	 = go
DEP    	 = dep
GOLINT 	 = golint
GREP   	 = grep
NPM    	 = npm
PRETTIER = prettier
NG       = ng
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
_make: deps build fe cleanup

PHONY += build
build: $(BIN) 

PHONY += deps
deps:
	$(DEP) ensure -v
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
	$(GO) test -v -cover ./...

PHONY += lint
lint:
	$(GOLINT) ./... | $(GREP) -v vendor || true

PHONY += run
run:
	$(GO) run -v \
		$(CURDIR)/cmd/$(APPNAME)/*.go -c $(CONFIG)

PHONY += cleanup
cleanup:

PHONY += fe
fe:
	cd $(CURDIR)/web && \
		$(NG) build --prod=true

PHONY += runfe
runfe:
	cd ./web && \
		$(NG) serve --port=8081

PHONY += prettify
prettify:
	$(PRETTIER) \
	    --config $(PRETTIER_CFG) \
	    --write \
	    	$(CURDIR)/web/src/**/*.js \
	    	$(CURDIR)/web/src/**/**/*.js \
	    	$(CURDIR)/web/src/**/*.vue \
	    	$(CURDIR)/web/src/**/**/*.vue

PHONY += help
help:
	@echo "Available targets:"
	@echo "  #        - creates binary in ./bin"
	@echo "  cleanup  - tidy up temporary stuff created by build or scripts"
	@echo "  deps     - ensure dependencies are installed"
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