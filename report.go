// report.go

/*
This file contains func for generating the report
*/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/atotto/encoding/csv"
	"github.com/ivpusic/grpool"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"
)

var (
	// this is to protect the Sites array
	lock sync.Mutex
)

// Private functions

// getResults read the JSON array generated and gone through jq
func getResults(file string) (res []byte, err error) {
	fh, err := os.Open(file)
	if err != nil {
		return res, errors.Wrapf(err, "can not open %s", file)
	}
	defer fh.Close()

	res, err = ioutil.ReadAll(fh)
	return res, errors.Wrapf(err, "can not read json %s", file)
}

func getSSLablsVersion(site ssllabs.Host) string {
	debug("%#v", site)
	return fmt.Sprintf("%s/%s", site.EngineVersion, site.CriteriaVersion)
}

// NewTLSReport generates everything we need for display/export
func NewTLSReport(reports []ssllabs.Host) (e *TLSReport, err error) {
	if len(reports) == 0 {
		return nil, fmt.Errorf("empty list")
	}

	e = &TLSReport{
		Date:    time.Now(),
		SSLLabs: getSSLablsVersion(reports[0]),
	}

	verbose("%d sites found.\n", len(reports))

	pool := grpool.NewPool(fJobs, len(reports))

	// release resources used by pool
	defer pool.Release()

	pool.WaitCount(len(reports))

	// Now analyze each site
	for _, site := range reports {
		debug("queueing %s\n", site.Host)

		current := site
		pool.JobQueue <- func() {
			completed := NewTLSSite(current)

			// Block on mutex
			lock.Lock()
			e.Sites = append(e.Sites, completed)
			lock.Unlock()

			pool.JobDone()
		}
	}

	pool.WaitAll()
	verbose("got all %d sites\n", len(e.Sites))
	debug("all=%v\n", e.Sites)
	sort.Sort(ByAlphabet(*e))
	return e, nil
}

type Types struct {
	Corrects map[string]int
	Insecure int
	ToFix    int
}

func (r *TLSReport) ColourMap(criteria string) Types {
	t := Types{Corrects: map[string]int{}}

	for _, site := range r.Sites {
		switch site.Type {
		case TypeHTTPSok:
			t.Corrects[selectColours(criteria)]++
		case TypeHTTPSnok:
			t.ToFix++
		case TypeHTTP:
			t.Insecure++
		}
	}
	return t
}

// ToCSV output a CSV file from a report
func (r *TLSReport) ToCSV(w io.Writer) (err error) {
	wh := csv.NewWriter(w)
	debug("%v\n", r.Sites)
	if err = wh.WriteStructHeader(r.Sites[0]); err != nil {
		return errors.Wrap(err, "can not write csv header")
	}

	err = wh.WriteStructAll(r.Sites)
	return errors.Wrap(err, "can not write csv file")
}

func (r *TLSReport) WriteCSV(w io.Writer, cntrs, https map[string]int) error {
	debug("WriteCSV")
	if len(r.Sites) == 0 {
		return fmt.Errorf("empty r")
	}

	if err := r.ToCSV(w); err != nil {
		return errors.Wrap(err, "Error can not generate CSV")
	}
	fmt.Fprintf(w, "\nTLS Summary\n")
	if err := writeSummary(os.Stdout, tlsKeys, cntrs); err != nil {
		fmt.Fprintf(os.Stderr, "can not generate TLS summary: %v", err)
	}
	fmt.Fprintf(w, "\nHTTP Summary\n")
	if err := writeSummary(os.Stdout, httpKeys, https); err != nil {
		fmt.Fprintf(os.Stderr, "can not generate HTTP summary: %v", err)
	}
	return nil
}
