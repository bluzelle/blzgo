KEY=$(key)

help:
	@echo "\
make read key=foo\
"

read:
	@go run examples/crud/read/*.go $(KEY)

hello_world:
	@go run examples/hello_world/*.go

pkgs:
	@dep ensure

test:
	@go test . -test.v

fmt:
	@gofmt -w *.go
	@gofmt -w examples/**/*.go

.PHONY: fmt test pkgs hello_world read help
