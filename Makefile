NAME=itunes-search-go

SHELL := env VERSION=$(VERSION) $(SHELL)
VERSION ?= $(shell git describe --tags $(git rev-list --tags --max-count=1))


# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


TARGET_MAX_CHAR_NUM=20


define colored
	@echo '${GREEN}$1${RESET}'
endef

## Show help
help:
	${call colored, help is running...}
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)


# Vet project
vet:
	${call colored, vet is running...}
	@go vet $(go list ./...)
.PHONY: vet

# Test project
test:
	${call colored, test is running...}
	@go test -race -v $(go list ./...)
.PHONY: test

# Sync vendor
sync-vendor:
	${call colored, sync-vendor is running...}
	@go mod tidy -v
	@go mod vendor
	@go mod verify
.PHONY: sync-vendor
