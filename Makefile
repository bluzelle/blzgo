KEY=$(key)
VALUE=$(value)
NEW_KEY=$(newkey)
UUID=$(uuid)
PROVE=$(prove)

help:
	@echo "\n\
make create key=foo value=bar\n\
make update key=foo value=bar\n\
make delete key=foo\n\
make rename key=foo newkey=baz\n\
\n\
make read key=foo prove=true\n\
make has key=foo\n\
make keys uuid=foo\n\
make keyvalues uuid=foo\n\
make count uuid=foo\n\
\n\
make txread key=foo prove=true\n\
make txhas key=foo\n\
make txkeys uuid=foo\n\
\n\
make account\n\
make version\n\
"

create:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

update:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

delete:
	@go run examples/crud/$@/*.go $(KEY)

rename:
	@go run examples/crud/$@/*.go $(KEY) $(NEW_KEY)

#

read:
	@go run examples/crud/$@/*.go $(KEY) $(PROVE)

has:
	@go run examples/crud/$@/*.go $(KEY)

keys:
	@go run examples/crud/$@/*.go $(UUID)

keyvalues:
	@go run examples/crud/$@/*.go $(UUID)

count:
	@go run examples/crud/$@/*.go $(UUID)

#

txread:
	@go run examples/crud/$@/*.go $(KEY) $(PROVE)

txhas:
	@go run examples/crud/$@/*.go $(KEY)

txkeys:
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

.PHONY: fmt test pkgs hello_world version account txkeys txhas txread count keyvalues keys has read rename delete update create help
