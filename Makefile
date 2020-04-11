hello_world:
	@go run samples/hello_world/*.go

pkgs:
	@dep ensure

test:
	@go test . -test.v

fmt:
	@gofmt -w *.go

.PHONY: fmt test pkgs hello_world
