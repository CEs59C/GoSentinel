.PHONY: clean build deploy run scp-env

BINARY_NAME = sentinel
BUILD_DIR = ./build
CMD_PATH = ./cmd/sentinel

ENV_FILE = .env.encrypted
KEY_FILE := .env.key
ENCRYPTION_KEY := $(shell cat $(BUILD_DIR)/$(KEY_FILE) 2>/dev/null)

#GOOS ?= linux
#GOARCH ?= amd64
#HOST ?= remote

GOOS ?= darwin
GOARCH ?= arm64
HOST ?= local

#GOOS ?= linux
#GOARCH ?= arm64
#HOST ?= local

BINARY = $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)
REMOTE_BINARY = $(BINARY_NAME)

ifeq ($(HOST),remote)
	HOST_ADDR = 85.155.101.203
	HOST_PORT = 22
	HOST_USER = root
	HOST_PATH = ~/
else ifeq ($(HOST),local)
	HOST_ADDR = 127.0.0.1
	HOST_PORT = 2222
	HOST_USER = q
	HOST_PATH = ~/
else
$(error Unknown HOST. Use HOST=remote or HOST=local)
endif

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)-*

encrypt:
	@echo "[Make] Шифрование .env..."
	go run tools/encrypt.go

build:
	@echo "[Make] Building for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY) $(CMD_PATH)


clean-remote:
	@echo "[Make] Cleaning to $(HOST)..."
	ssh $(HOST_USER)@$(HOST_ADDR) -p $(HOST_PORT) "rm -f ~/$(BINARY_NAME)*"
	ssh $(HOST_USER)@$(HOST_ADDR) -p $(HOST_PORT) "rm -f ~/$(ENV_FILE)"

scp-env:
	@echo "[Make] Uploading env file..."
	scp -P $(HOST_PORT) $(BUILD_DIR)/$(ENV_FILE) \
	$(HOST_USER)@$(HOST_ADDR):$(HOST_PATH)

deploy: clean-remote build scp-env
	@echo "[Make] Deploying to $(HOST)..."
	scp -P $(HOST_PORT) $(BINARY) \
	$(HOST_USER)@$(HOST_ADDR):$(HOST_PATH)/$(REMOTE_BINARY)

run: deploy
	@echo "[Make] Running on $(HOST)..."
	ssh $(HOST_USER)@$(HOST_ADDR) -p $(HOST_PORT) "ENCRYPTION_KEY='$(ENCRYPTION_KEY)' ~/$(REMOTE_BINARY)" > local_result.txt 2>&1
	@cat local_result.txt

run-local: build
	@echo "[Make] Local run for $(GOOS)/$(GOARCH)..."
	@if [ -f $(BUILD_DIR)/$(KEY_FILE) ]; then \
		echo "[Make] Ключ найден, запускаем с шифрованием"; \
		ENCRYPTION_KEY="$(ENCRYPTION_KEY)" $(BINARY); \
	else \
		echo "[Make] Ключ не найден, работаем с нешифрованным .env"; \
		$(BINARY); \
	fi

run-remote-amd64:
	$(MAKE) run GOOS=linux GOARCH=amd64 HOST=remote

run-local-arm64:
	$(MAKE) run GOOS=linux GOARCH=arm64 HOST=local

run-mac:
	$(MAKE) run-local GOOS=darwin GOARCH=arm64 HOST=local

install:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
lint:
	@golangci-lint run ./... -v
lint_autofix:
	@G0111MODULE=on $(GOLINT) run ./... -v --fix