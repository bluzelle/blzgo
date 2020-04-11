KEY=$(key)
VALUE=$(value)

help:
	@echo "\
make read key=foo\
make create key=foo value=bar\
make update key=foo value=bar\
make account\
"

read:
	@go run examples/crud/$@/*.go $(KEY)

create:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

update:
	@go run examples/crud/$@/*.go $(KEY) $(VALUE)

account:
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

.PHONY: fmt test pkgs hello_world account update create read help
