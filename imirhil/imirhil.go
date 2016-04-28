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
	baseURL   = "https://tls.imirhil.fr/https/"
	ext       = ".json"
	aDay      = time.Duration(24) * time.Hour
	threshold = time.Duration(30) * aDay
)

// Private area

// callAPI makes the actual call, probably "Pending" as 1st answer
func callAPI(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusOK {
		var body []byte

		body, err = ioutil.ReadAll(resp.Body)
		if string(body) == "pending" {
			resp, err = http.Get(url)
			if err != nil {
				return
			}
		}
	}
	return
}

// Public functions

// GetScore retrieves the current score for tls.imirhil.fr
func GetScore(site string) (score string) {
	full, err := GetDetailedReport(site)
	if err != nil {
		score = "Z"
		log.Printf("Error: can not get imirhil rating: %v", err)
		return
	}
	score = full.Score.Rank
	return
}

// GetDetailedReport retrieve the full data
func GetDetailedReport(site string) (report Report, err error) {
	var body []byte

	resp, err := http.Get(baseURL + site + ext)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusOK {

		body, err = ioutil.ReadAll(resp.Body)
		if string(body) == "pending" {
			resp, err = http.Get(baseURL + site + ext)
			if err != nil {
				return
			}
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("did not get acceptable status code: %v body: %q", resp.Status, body)
		return
	}

	err = json.Unmarshal(body, &report)
	return
}