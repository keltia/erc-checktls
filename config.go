package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	dbrcFile = filepath.Join(os.Getenv("HOME"), ".dbrc")

	user, password string
)

func setupProxyAuth(ctx *Context) (err error) {
	err = loadDbrc(ctx, dbrcFile)
	if err != nil {
		log.Printf("No dbrc file: %v", err)
	}
	if fVerbose {
		log.Printf("Proxy user %s found.", user)
	}

	// Do we have a proxy user/password?
	if user != "" && password != "" {
		auth := fmt.Sprintf("%s:%s", user, password)
		ctx.proxyauth = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	}

	return
}

func loadDbrc(ctx *Context, file string) (err error) {
	fh, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Error: can not find %s: %v", file, err)
	}
	defer fh.Close()

	/*
	   Format:
	   <db>     <user>    <pass>   <type>
	*/
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		// Replace all tabs by a single space
		l := strings.Replace(line, "\t", " ", -1)
		flds := strings.Split(l, " ")

		// Check what we need
		if flds[0] == "proxy" {
			user = flds[1]
			password = flds[2]
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading dbrc %s", dbrcFile)
	}

	if user == "" {
		return fmt.Errorf("no user/password for cimbl in %s", dbrcFile)
	}

	return
}
