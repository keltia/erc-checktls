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
	baseURL = "https://tls.imirhil.fr/https/"
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
		ctx.fVerbose = true
	}
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

	str := baseURL + site + ext
	req, trsp := setupTransport(str)

	if req == nil || trsp == nil {
		err = fmt.Errorf("Can not setup connection")
		return
	}

	// It is better to re-use than creating a new one each time
	if ctx.Client == nil {
		ctx.Client = &http.Client{Transport: trsp, Timeout: DefaultWait}
	}

	resp, err := ctx.Client.Do(req)
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
