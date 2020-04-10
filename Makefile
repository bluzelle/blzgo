sample:
	@go run samples/*.go

pkgs:
	@dep ensure

test:
	@go test ./src -test.v

fmt:
	@gofmt -w samples/*
	@gofmt -w src/*

.PHONY: fmt test pkgs sample
