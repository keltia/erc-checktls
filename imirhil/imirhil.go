// imirhil.go

/*
  This file contains the datatypes used by tls.imirhil.fr
*/

package imirhil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "https://tls.imirhil.fr/"
	typeURL = "https/"
	ext     = ".json"

	DefaultWait = 10 * time.Second
	Version     = "201712"
)

var (
	ctx = &Context{}
)

// Private area

// Public functions

// Init setups proxy authentication
func Init(fVerbose bool, proxyauth string) {
	if proxyauth != "" {
		ctx.proxyauth = proxyauth
	}
	if fVerbose {
		ctx.verbose = true
	}

	_, trsp := setupTransport(baseURL)
	ctx.Client = &http.Client{Transport: trsp, Timeout: DefaultWait}
	verbose("imirhil: ctx=%#v", ctx)
}

// GetScore retrieves the current score for tls.imirhil.fr
func GetScore(site string) (score string) {
	full, err := GetDetailedReport(site)
	if err != nil {
		score = "Z"
		log.Printf("Error: can not get imirhil rating: %v", err)
		return
	}
	score = full.Hosts[0].Grade.Rank
	return
}

// GetDetailedReport retrieve the full data
func GetDetailedReport(site string) (report Report, err error) {
	var body []byte

	str := fmt.Sprintf("%s/%s/%s.%s", baseURL, typeURL, site, ext)

	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		log.Printf("error: req is nil: %v", err)
		return Report{}, nil
	}

	verbose("req=%#v", req)
	verbose("clt=%#v", ctx.Client)

	resp, err := ctx.Client.Do(req)
	if err != nil {
		verbose("err=%#v", err)
		return
	}
	verbose("resp=%#v", resp)
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {

		if string(body) == "pending" {
			time.Sleep(10 * time.Second)
			resp, err = ctx.Client.Do(req)
			if err != nil {
				return
			}
		}
	} else {
		err = fmt.Errorf("did not get acceptable status code: %v body: %q", resp.Status, body)
		return
	}

	err = json.Unmarshal(body, &report)
	return
}
