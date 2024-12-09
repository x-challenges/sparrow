YDB_HOST=ydb.serverless.yandexcloud.net:2135
YDB_DATABASE=/ru-central1/b1gcm11knnnopur9fil6/etn4hff7q981ib0o7mnl
YDB_TOKEN=$(shell yc iam create-token)

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

migrate: .check-var-command
	goose -dir migrations ydb \
		"grpcs://$(YDB_HOST)/?database=$(YDB_DATABASE)&token=$(YDB_TOKEN)&go_query_mode=scripting&go_fake_tx=scripting&go_query_bind=declare,numeric" \
		$(command) -v
