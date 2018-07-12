
# Main Makefile for labs2pg
#
# Copyright 2015 © by Ollivier Robert for the EEC
#

.PATH= ssllabs
GOBIN=   ${GOPATH}/bin

SRCS= main.go cli.go report.go utils.go \
	config.go \
    ssllabs/ssllabs.go ssllabs/types.go

OPTS=	-ldflags="-s -w" -v

all: erc-checktls

erc-checktls: ${SRCS}
	go build ${OPTS}
	go test -v

install:
	go install ${OPTS}

lint:
	gometalinter .

clean:
	go clean -v

push:
	git push --all
	git push --tags
	git push --all backup
	git push --tags backup
	git push --all upstream
	git push --tags upstream
