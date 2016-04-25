# Main Makefile for labs2pg
#
# Copyright 2015 © by Ollivier Robert for the EEC
#

.PATH= imirhil
GOBIN=   ${GOPATH}/bin

all: erc-checktls

erc-checktls: labs2pg.go types.go results.go cli.go ssllabs.go imirhil/imirhil.go
	go build -v
	go test -v

install:
	go install -v

clean:
	go clean -v

push:
	git push --all
	git push --tags
	git push --all backup
	git push --tags backup
	git push --all upstream
	git push --tags upstream
