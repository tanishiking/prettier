GO=go

.PHONY: test
test:
	${GO} test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	go clean
