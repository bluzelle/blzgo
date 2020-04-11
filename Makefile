KEY=$(key)
VALUE=$(value)
NEW_KEY=$(newkey)
UUID=$(uuid)
PROVE=$(prove)

help:
	@echo "\n\
make read key=foo prove=true\n\
make create key=foo value=bar\n\
make update key=foo value=bar\n\
make delete key=foo\n\
make rename key=foo newkey=baz\n\
make has key=foo\n\
make keys uuid=foo\n\
make keyvalues uuid=foo\n\
make count uuid=foo\n\
\n\
make account\n\
make version\n\
"

read:
	@go run examples/crud/$@/*.go $(KEY) $(PROVE)

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

count:
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

.PHONY: fmt test pkgs hello_world version account count keyvalues keys has rename delete update create read help
