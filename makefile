COVER_FILE			?= coverage.out
SOURCE_PATHS		?= ./pkg/...

unit_test: ## Run unit tests
	go test -gcflags=all=-l -timeout=20m `go list $(SOURCE_PATHS)` ${TEST_FLAGS} -v

cover:  ## Generates coverage report
	go test -gcflags=all=-l -timeout=20m `go list $(SOURCE_PATHS)` -coverprofile $(COVER_FILE) ${TEST_FLAGS} -v
