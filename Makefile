LOG_LEVEL ?= info

.PHONY: build
build:
	go build -o dist/ ./...

.PHONY: test
test:
	@LOG_LEVEL=${LOG_LEVEL} go test -v ./...
