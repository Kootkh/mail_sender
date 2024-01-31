.PHONY: build
build:
	go build -v ./cmd/mail_sender

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build