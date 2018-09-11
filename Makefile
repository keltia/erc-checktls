
# Main Makefile for labs2pg
#
# Copyright 2015-2018 © by Ollivier Robert for the EEC
#

GO=		go
GOBIN=   ${GOPATH}/bin

SRCS= main.go categories.go cli.go html.go  report.go site.go utils.go types.go \
	main-packr.go

OPTS=	-ldflags="-s -w" -v

all: erc-checktls

main-packr.go: main.go
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
