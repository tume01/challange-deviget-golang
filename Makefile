.PHONY:fmt lint test

PACKAGES = $(shell go list ./...)

fmt:
	@echo "===> Formatting"
	@go fmt $(PACKAGES)

lint:
	@echo "===> Linting with vet"
	@go vet $(PACKAGES)

test: lint
	@echo "===> Testing"
	@go test -race -count=1 -covermode=atomic $(PACKAGES)