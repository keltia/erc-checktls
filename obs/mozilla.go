package obs

/*
Not going to implement the full scan report struct, I do not need it, juste grade/score
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keltia/proxy"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "https://http-observatory.security.mozilla.org/api/v1"

	// DefaultWait is the timeout
	DefaultWait = 10 * time.Second

	// MyVersion is the API version
	MyVersion = "0.2.0"

	// MyName is the name used for the configuration
	MyName = "obs"
)

// Private area

func myRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

func (c *Client) callAPI(site, word, sbody string) (*Analyze, error) {

	str := fmt.Sprintf("%s/analyze?host=%s", c.baseurl, site)

	c.debug("str=%s", str)
	req, err := http.NewRequest(word, str, nil)
	if err != nil {
		log.Printf("error: req is nil: %v", err)
		return &Analyze{}, errors.Wrap(err, "req is nil")
	}

	c.debug("req=%#v", req)
	c.debug("clt=%#v", c.client)

	// If we have a POST and a body, insert them.
	if sbody != "" && word == "POST" {
		body := []byte(sbody)
		buf := bytes.NewReader(body)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(buf.Len())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.verbose("err=%#v", err)
		return &Analyze{}, errors.Wrap(err, "1st call failed")
	}
	c.debug("resp=%#v", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Analyze{}, errors.Wrap(err, "can not read body")
	}

	if resp.StatusCode == http.StatusOK {

		c.debug("status OK")

		if string(body) == "pending" {
			time.Sleep(10 * time.Second)
			resp, err = c.client.Do(req)
			if err != nil {
				return &Analyze{}, errors.Wrap(err, "pending failed")
			}
			c.verbose("resp was %v", resp)
		}
	} else if resp.StatusCode == http.StatusFound {
		str := resp.Header["Location"][0]

		c.debug("Got 302 to %s", str)

		req, err = http.NewRequest("GET", str, nil)
		if err != nil {
			return &Analyze{}, errors.Wrap(err, "Cannot handle redirect")
		}

		resp, err = c.client.Do(req)
		if err != nil {
			return &Analyze{}, errors.Wrap(err, "client.Do failed")
		}
		c.verbose("resp was %v", resp)
	} else {
		return &Analyze{}, errors.Wrapf(err, "did not get acceptable status code: %v body: %q", resp.Status, body)
	}

	var report Analyze

	err = json.Unmarshal(body, &report)

	// Give some time to performe the test
	if report.State == "PENDING" {
		time.Sleep(2 * time.Second)
		err = nil
	}

	return &report, err
}

// Public functions

// NewClient setups proxy authentication
func NewClient(cnf ...Config) *Client {
	var c *Client

	// Set default
	if len(cnf) == 0 {
		c = &Client{
			baseurl: baseURL,
			timeout: DefaultWait,
		}
	} else {
		c = &Client{
			baseurl: cnf[0].BaseURL,
			level:   cnf[0].Log,
			refresh: cnf[0].Refresh,
		}

		if cnf[0].Timeout == 0 {
			c.timeout = DefaultWait
		} else {
			c.timeout = time.Duration(cnf[0].Timeout) * time.Second
		}

		// Ensure we have the API endpoint right
		if c.baseurl == "" {
			c.baseurl = baseURL
		}

		c.verbose("got cnf: %#v", cnf[0])
	}

	proxyauth, err := proxy.SetupProxyAuth()
	if err != nil {
		c.proxyauth = proxyauth
	}

	_, trsp := proxy.SetupTransport(c.baseurl)
	c.client = &http.Client{
		Transport:     trsp,
		Timeout:       c.timeout,
		CheckRedirect: myRedirect,
	}
	c.debug("mozilla: c=%#v", c)
	return c
}

// GetScore returns the integer value of the grade
func (c *Client) GetScore(site string) (score int, err error) {
	c.debug("GetScore")
	_, err = c.callAPI(site, "POST", "hidden=true&rescan=true")
	if err != nil {
		return -1, errors.Wrap(err, "callAPI failed")
	}
	r, err := c.callAPI(site, "GET", "")
	return r.Score, errors.Wrap(err, "GetScore failed")
}

// GetGrade returns the letter equivalent to the score
func (c *Client) GetGrade(site string) (grade string, err error) {
	c.debug("GetGrade")
	_, err = c.callAPI(site, "POST", "hidden=true&rescan=true")
	if err != nil {
		return "Z", errors.Wrap(err, "callAPI failed")
	}
	r, err := c.callAPI(site, "GET", "")
	return r.Grade, errors.Wrap(err, "GetGrade failed")
}
