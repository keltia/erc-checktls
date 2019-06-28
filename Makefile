
# Main Makefile for labs2pg
#
# Copyright 2015-2018 © by Ollivier Robert for the EEC
#

GO=		go
GOBIN=   ${GOPATH}/bin

NAME=	erc-checktls
SRCS= cmd/${NAME}/main.go cmd/${NAME}/cli.go html.go html_subr.go report.go \
	resources.go summaries.go utils.go types.go \
	TLS-packr.go \
	site/site.go site/utils.go site/types.go

OPTS=	-ldflags="-s -w" -v

all: ${NAME}

TLS-packr.go: report.go files/templ.html files/summaries.html files/sites-list.csv
	packr2

${NAME}: ${SRCS}
	${GO} build ${OPTS} ./cmd/...

install: all test
	${GO} install ${OPTS} ./cmd/...

test: all
	${GO} test ./...

clean:
	${GO} clean -v
	packr2 clean

push:
	git push --all
	git push --tags
	git push --all backup
	git push --tags backup
