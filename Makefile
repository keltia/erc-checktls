
# Main Makefile for labs2pg
#
# Copyright 2015-2018 © by Ollivier Robert for the EEC
#

GO=		go
GOBIN=   ${GOPATH}/bin

SRCS= main.go categories.go cli.go html.go html_subr.go report.go \
	resources.go site.go summaries.go utils.go types.go \
	main-packr.go

OPTS=	-ldflags="-s -w" -v

all: erc-checktls

main-packr.go: main.go files/templ.html files/summaries.html files/sites-list.csv
	packr2

erc-checktls: ${SRCS}
	${GO} build ${OPTS}

install: all test
	${GO} install ${OPTS}

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
