KEY=$(key)
VALUE=$(value)
NEW_KEY=$(newkey)

help:
	@echo "\
make read key=foo\
make create key=foo value=bar\
make update key=foo value=bar\
make delete key=foo\
make rename key=foo newkey=baz\
make account\
make version\
"

read:
	@go run examples/crud/$@/*.go $(KEY)

create:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

update:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

delete:
	@go run examples/crud/$@/*.go $(KEY)

rename:
	@go run examples/crud/$@/*.go $(KEY) $(NEW_KEY)

account:
	@go run examples/crud/$@/*.go

version:
	@go run examples/crud/$@/*.go

hello_world:
	@go run examples/hello_world/*.go

pkgs:
	@dep ensure

test:
	@go test . -test.v

fmt:
	@gofmt -w *.go
	@gofmt -w examples/**/*.go

.PHONY: fmt test pkgs hello_world version account rename delete update create read help
