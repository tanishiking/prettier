GO=go

.PHONY: test
test:
	${GO} test -v ./...

.PHONY: lintdeps
lintdeps:
	GO111MODULE=off go get golang.org/x/lint/golint

.PHONY: lint
lint: lintdeps
	go vet
	golint -set_exit_status

.PHONY: clean
clean:
	go clean
