KEY=$(key)
VALUE=$(value)
NEW_KEY=$(newkey)
PROVE=$(prove)
LEASE=$(lease)
N=$(n)
test=$(o)

help:
	@echo "\n\
make create key=foo value=bar lease=0\n\
make update key=foo value=bar lease=0\n\
make delete key=foo\n\
make rename key=foo newkey=baz\n\
make deleteall\n\
make multiupdate\n\
make renewlease key=foo lease=0\n\
make renewallleases lease=0\n\
\n\
make read key=foo prove=true\n\
make has key=foo\n\
make keys\n\
make keyvalues\n\
make count\n\
make getlease key=foo\n\
make getnshortestleases n=2\n\
\n\
make txread key=foo prove=true\n\
make txhas key=foo\n\
make txkeys\n\
make txkeyvalues\n\
make txcount\n\
make txgetlease key=foo\n\
make txgetnshortestleases n=2\n\
\n\
make account\n\
make version\n\
\n\
make hello_world\n\
make multi\n\
make uuid\n\
\n\
make test\n\
"

create:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE) $(LEASE)

update:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE) $(LEASE)

delete:
	@go run examples/crud/$@/*.go $(KEY)

rename:
	@go run examples/crud/$@/*.go $(KEY) $(NEW_KEY)

deleteall:
	@go run examples/crud/$@/*.go

multiupdate:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

renewlease:
	@go run examples/crud/$@/*.go $(KEY) $(LEASE)

renewallleases:
	@go run examples/crud/$@/*.go $(LEASE)

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

getnshortestleases:
	@go run examples/crud/$@/*.go $(N)

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

txgetnshortestleases:
	@go run examples/crud/$@/*.go $(N)

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
	@go test . -test.v
#	@./test.sh

test-spec:
	@go test . -test.v -run 'Test$(o)'

fmt:
	@gofmt -w *.go
	@gofmt -w examples/**
	@gofmt -w uat/**

uat: fmt
	@go run uat/*.go

.PHONY: help \
	create \
	update \
	delete \
	rename \
	deleteall \
	multiupdate \
	renewlease \
	renewallleases \
	read \
	has \
	keys \
	keyvalues \
	count \
	getlease \
	getnshortestleases \
	txread \
	txhas \
	txkeys \
	txkeyvalues \
	txcount \
	txgetlease \
	txgetnshortestleases \
	account \
	version \
	hello_world \
	multi \
	uuid \
	pkgs \
	test \
	test-spec \
	fmt \
	uat
