# Main Makefile for labs2pg
#
# Copyright 2015 © by Ollivier Robert for the EEC
#

GOBIN=   ${GOPATH}/bin

all: labs2pg.go types.go results.go cli.go
	go build -v ./...
	go test -v ./...

install:
	go install -v

clean:
	go clean -v

push:
	git push --all
	git push --all upstream
	git push --all backup
	git push --tags
	git push --tags upstream
	git push --tags backup
