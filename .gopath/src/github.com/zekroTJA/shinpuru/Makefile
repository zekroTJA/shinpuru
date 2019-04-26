GO  	= go
DEP 	= dep

PACKAGE = github.com/zekroTJA/shinpuru
GOPATH  = $(CURDIR)/.gopath
WDIR    = $(GOPATH)/src/$(PACKAGE)

BINNAME = shinpuru
BINLOC  = $(CURDIR)

ifeq ($(OS),Windows_NT)
	EXTENSION=.exe
endif

BIN = $(BINLOC)/$(BINNAME)$(EXTENSION)

TAG        = $(shell git describe --tags)
COMMIT     = $(shell git rev-parse HEAD)

SQLLDFLAGS = $(shell bash ./scripts/getsqlschemes.bash)

.PHONY: _make installdeps cleanup _finish run

_make: $(WDIR) $(BIN) cleanup _finish

$(WDIR):
	@echo [ INFO ] creating working directory '$@'...
	mkdir -p $@
	cp -R $(CURDIR)/* $@/ 

$(BIN): installdeps
	@echo [ INFO ] building binary '$(BIN)'...
	(env GOPATH=$(GOPATH) $(GO) build -v -o $@ -ldflags "\
		-X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$(TAG) \
		-X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$(COMMIT) \
		-X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE \
		$(SQLLDFLAGS)" \
		$(WDIR)/cmd/shinpuru)

installdeps:
	@echo [ INFO ] installing dependencies with dep...
	cd $(WDIR) && \
		$(DEP) ensure

cleanup:
	@echo [ INFO ] cleaning up...
	rm -r -f $(GOPATH)

_finish:
	@echo ------------------------------------------------------------------------------
	@echo [ INFO ] Build successful.
	@echo [ INFO ] Your build is located at '$(BIN)'

run:
	(env GOPATH=$(CURDIR)/../../../.. $(GO) run -v ./cmd/shinpuru -c ./config/private.config.yaml)