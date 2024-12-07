GOCMD=go

.PHONY: lint test cover generate

.check-var-%:
	@[ "$($*)" ] || (echo "$* is undefined"; exit 1)

.yesno:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

lint:
	golangci-lint run --fix

test:
	${GOCMD} test -race -v -coverpkg=./... -coverprofile=coverage.out ./...

cover: test
	${GOCMD} tool cover -func=coverage.out | tail -1

generate:
	${GOCMD} generate ./...
