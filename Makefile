APP := tweet-me
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)
GO := go

.PHONY: all build install uninstall clean release version

all: build

build:
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(APP) .

install: build
	mkdir -p $(HOME)/.local/bin
	cp $(APP) $(HOME)/.local/bin/$(APP)
	@echo "Installed $(APP) to $$HOME/.local/bin (ensure it's on your PATH)"

uninstall:
	rm -f $(HOME)/.local/bin/$(APP)
	@echo "Removed $(APP) from $$HOME/.local/bin"

clean:
	rm -f $(APP)
	@echo "Cleaned build artifacts"

version:
	@echo $(VERSION)

release:
	GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(APP)-$(VERSION)-linux-amd64
	GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(APP)-$(VERSION)-darwin-arm64
	GOOS=darwin GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(APP)-$(VERSION)-darwin-amd64
	GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(APP)-$(VERSION)-windows-amd64.exe
	@echo "Artifacts in dist/"
