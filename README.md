erc-checktls
============

[![GitHub release](https://img.shields.io/github/release/keltia/erc-checktls.svg)](https://github.com/keltia/erc-checktls/releases) 
[![GitHub issues](https://img.shields.io/github/issues/keltia/erc-checktls.svg)](https://github.com/keltia/erc-checktls/issues)
[![Go Version](https://img.shields.io/badge/go-1.10-blue.svg)](https://golang.org/dl/)
[![Build Status](https://travis-ci.org/keltia/erc-checktls.svg?branch=master)](https://travis-ci.org/keltia/erc-checktls)
[![GoDoc](http://godoc.org/github.com/keltia/erc-checktls?status.svg)](http://godoc.org/github.com/keltia/erc-checktls)
[![SemVer](http://img.shields.io/SemVer/2.0.0.png)](https://semver.org/spec/v2.0.0.html)
[![License](https://img.shields.io/pypi/l/Django.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/keltia/erc-checktls)](https://goreportcard.com/report/github.com/keltia/erc-checktls)

This is a small utility which will provide summary & diff-like operations for the reports generated by [ssllabs-scan](https://github.com/ssllabs/ssllabs-scan).

In addition the grade checked by [Imirhil](https://tls.imirhil.fr/) will be checked as well and displayed.  We now retrieve the [Mozilla Observatory](https://observatory.mozilla.org/) grade as well.

## Requirements

* Go >= 1.10
* jq (optional)

You need to install three of my modules if you are using Go 1.10.x or earlier.

    go get github.com/keltia/proxy
    go get github.com/keltia/cryptcheck
    go get github.com/keltia/observatory

I also use a number of external modules:

	github.com/atotto/encoding/csv
	github.com/gobuffalo/packr
	github.com/ivpusic/grpool
	github.com/pkg/errors
	github.com/olekukonko/tablewriter

If you want to run `make test` you will need these:

	github.com/stretchr/testify/assert
	github.com/stretchr/testify/require

With Go 1.11+ and its modules support, it should work out of the box with

    go get github.com/keltia/erc-checktls

if you have `GO111MODULE` set to `on`.

## Usage

SYNOPSIS
```
erc-checktls [-vDIMV] [-j N] [-t csv|text|html] [-o file] [-s file] [-S site] <json file>
  
  -D	Debug mode
  -I	Do not fetch tls.imirhil.fr grade
  -M	Do not fetch Mozilla Observatory data
  -R	Force refresh
  -S string
    	Display that site
  -j    Set the # of parallel jobs to run (default # of cores)
  -o string
    	Save into file (default stdout) (default "-")
  -s string
    	Save summary there (default "summary")
  -t string
    	Type of report (default "csv")
  -v	Verbose mode
  
If you just want to find all wildcard certificates use this:

  -wild
    	Display wildcards
```

The json file needs to be generated by running `ssllabs-scan` post-processed by `jq` like this:
 
```
ssllabs-scan -hostfile <host list> > <json file>
```

OPTIONS

| Option  | Default | Description|
| ------- |---------|------------|
| -D      | false   | Debug mode |
| -I      | false   | Do not fetch tls.imirhil.fr grade |
| -M      | false   | Do not fetch Mozilla Observatory data |
| -R      | false   | Force refresh |
| -S      | none    | Displays that site info only |
| -j      | # cores | Set level of parallelism (default # of CPU cores |
| -o      | -       | Output into that file (default stdout) |
| -s      | summary | Save summary in that file (default summary.html) |
| -t      | csv     | Output plain text, html or csv |
| -v      | false   | Be verbose |
| -wild   |         | Report wildcard certificates |

## Using behind a web Proxy

Dependency: proxy support is provided by my `github.com/keltia/proxy` module.

UNIX/Linux:

```
    export HTTP_PROXY=[http://]host[:port] (sh/bash/zsh)
    setenv HTTP_PROXY [http://]host[:port] (csh/tcsh)
```

Windows:

```
    set HTTP_PROXY=[http://]host[:port]
```

The rules of Go's `ProxyFromEnvironment` apply (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, lowercase variants allowed).

If your proxy requires you to authenticate, please create a file named `.netrc` in your HOME directory with permissions either `0400` or `0600` with the following data:

    machine proxy user <username> password <password>
    
and it should be picked up. On Windows, the file will be located at

    %LOCALAPPDATA%\proxy\netrc

## TODO

- Implement full online calls for SSLLabs
- Better separation between batch & online mode

## License

The [BSD 2-Clause license](https://github.com/keltia/erc-checktls/blob/master/LICENSE).

# Contributing

This project is an open Open Source project, please read `CONTRIBUTING.md`.

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.
