LOG_LEVEL ?= info

.PHONY: build
build:
	go build -o dist/ ./...

.PHONY: test
test:
	@LOG_LEVEL=${LOG_LEVEL} go test -v ./... ${TESTFLAGS}

.PHONY: mock
mock:
	mockgen -source=domain/types.go -destination=./domain/mock/types_mock.go -package=domain
	mockgen -source=application/types.go -destination=./application/mock/types_mock.go -package=application

.PHONY: run
run:
	go run main.go
