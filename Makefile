KEY=$(key)
VALUE=$(value)
NEW_KEY=$(newkey)
UUID=$(uuid)

help:
	@echo "\
make read key=foo\
make create key=foo value=bar\
make update key=foo value=bar\
make delete key=foo\
make rename key=foo newkey=baz\
make has key=foo\
make keys uuid=foo\
make keyvalues uuid=foo\
\
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

has:
	@go run examples/crud/$@/*.go $(KEY)

keys:
	@go run examples/crud/$@/*.go $(UUID)

keyvalues:
	@go run examples/crud/$@/*.go $(UUID)

#

account:
	@go run examples/crud/$@/*.go

version:
	@go run examples/crud/$@/*.go

#

hello_world:
	@go run examples/hello_world/*.go

pkgs:
	@dep ensure

test:
	@go test . -test.v

fmt:
	@gofmt -w *.go
	@gofmt -w examples/**/*.go

.PHONY: fmt test pkgs hello_world version account keyvalues keys has rename delete update create read help
