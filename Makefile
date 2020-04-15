KEY=$(key)
VALUE=$(value)
NEW_KEY=$(newkey)
PROVE=$(prove)
LEASE=$(lease)

help:
	@echo "\n\
make create key=foo value=bar lease=0\n\
make update key=foo value=bar lease=0\n\
make delete key=foo\n\
make rename key=foo newkey=baz\n\
\n\
make read key=foo prove=true\n\
make has key=foo\n\
make keys\n\
make keyvalues\n\
make count\n\
make getlease key=foo\n\
\n\
make txread key=foo prove=true\n\
make txhas key=foo\n\
make txkeys\n\
make txkeyvalues\n\
make txcount\n\
make txgetlease key=foo\n\
\n\
make account\n\
make version\n\
"

create:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE) $(LEASE)

update:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE) $(LEASE)

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
	@go run examples/crud/$@/*.go

keyvalues:
	@go run examples/crud/$@/*.go

count:
	@go run examples/crud/$@/*.go

getlease:
	@go run examples/crud/$@/*.go $(KEY)

#

txread:
	@go run examples/crud/$@/*.go $(KEY) $(PROVE)

txhas:
	@go run examples/crud/$@/*.go $(KEY)

txkeys:
	@go run examples/crud/$@/*.go

txkeyvalues:
	@go run examples/crud/$@/*.go

txcount:
	@go run examples/crud/$@/*.go

txgetlease:
	@go run examples/crud/$@/*.go $(KEY)

#

deleteall:
	@go run examples/crud/$@/*.go

multiupdate:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

#

account:
	@go run examples/crud/$@/*.go

version:
	@go run examples/crud/$@/*.go

#

hello_world:
	@go run examples/$@/*.go

multi:
	@go run examples/$@/*.go

uuid:
	@go run examples/$@/*.go

pkgs:
	@dep ensure

test:
#	@go test . -test.v
	@./test.sh

fmt:
	@gofmt -w *.go
	@gofmt -w examples/**/*.go

.PHONY: fmt \
	test \
	pkgs \
	uuid \
	multi \
	hello_world \
	version \
	account \
	multiupdate \
	deleteall \
	txgetlease \
	txcount \
	txkeyvalues \
	txkeys \
	txhas \
	txread \
	getlease \
	count \
	keyvalues \
	keys \
	has \
	read \
	rename \
	delete \
	update \
	create \
	help
