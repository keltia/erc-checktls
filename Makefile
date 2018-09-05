
# Main Makefile for labs2pg
#
# Copyright 2015-2018 © by Ollivier Robert for the EEC
#

GO=		go
.PATH= ssllabs
GOBIN=   ${GOPATH}/bin

SRCS= main.go cli.go report.go utils.go \
	main-packr.go

OPTS=	-ldflags="-s -w" -v

all: erc-checktls

main-packr.go:
	packr

erc-checktls: ${SRCS}
	${GO} build ${OPTS}
	${GO} test -v

install:
	${GO} install ${OPTS}

lint:
	gometalinter .

clean:
	${GO} clean -v
	packr clean

push:
	git push --all
	git push --tags
	git push --all backup
	git push --tags backup
