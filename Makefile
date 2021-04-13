PACKAGES = $(shell go list ./...)

lint:
	@echo "===> Linting with vet"
	@go vet $(PACKAGES)

test: lint
	@echo "===> Testing"
	@go test -race -count=1 -covermode=atomic $(PACKAGES)