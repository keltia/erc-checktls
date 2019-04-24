// report.go

/*
This file contains func for generating the report
*/
package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/atotto/encoding/csv"
	"github.com/ivpusic/grpool"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"

	"github.com/keltia/erc-checktls/site"
)

// NewTLSReport generates everything we need for display/export
func NewTLSReport(reports []ssllabs.Host) (r *TLSReport, err error) {
	// this is to protect the Sites array
	var lock sync.Mutex

	if len(reports) == 0 {
		return nil, fmt.Errorf("empty list")
	}

	r = &TLSReport{
		Date:    time.Now(),
		SSLLabs: getSSLablsVersion(reports[0]),
		cntrs:   map[string]int{},
		https:   map[string]int{},
	}

	verbose("%d sites found.\n", len(reports))

	pool := grpool.NewPool(fJobs, len(reports))

	// release resources used by pool
	defer pool.Release()

	pool.WaitCount(len(reports))

	// Setup our env.
	site.Init(site.Flags{
		LogLevel:      logLevel,
		IgnoreMozilla: fIgnoreMozilla,
		IgnoreImirhil: fIgnoreImirhil,
		Contracts:     contracts,
	})

	// Now analyze each site
	for _, s := range reports {
		debug("queueing %s\n", s.Host)

		current := s
		pool.JobQueue <- func() {
			// Block on mutex
			lock.Lock()
			s := site.NewFromHost(current)

			r.GatherStats(s)
			r.Sites = append(r.Sites, s)
			lock.Unlock()

			pool.JobDone()
		}
	}

	pool.WaitAll()
	verbose("got all %d sites\n", len(r.Sites))
	debug("all=%v\n", r.Sites)
	sort.Sort(ByAlphabet(*r))
	return r, nil
}

type Types struct {
	Corrects map[string]int
	Insecure int
	ToFix    int
}

func (r *TLSReport) ColourMap(criteria string) Types {
	t := Types{Corrects: map[string]int{}}

	for _, s := range r.Sites {
		switch s.Type {
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

func (r *TLSReport) WriteCSV(w io.Writer) error {
	debug("WriteCSV\n")
	debug("r=%#v\n", r)
	if len(r.Sites) == 0 {
		return fmt.Errorf("empty r")
	}

	if err := r.ToCSV(w); err != nil {
		return errors.Wrap(err, "Error can not generate CSV")
	}
	fmt.Fprintf(w, "\nTLS Summary\n")
	if err := writeSummary(w, tlsKeys, r.cntrs); err != nil {
		fmt.Fprintf(os.Stderr, "can not generate TLS summary: %v", err)
	}
	fmt.Fprintf(w, "\nHTTP Summary\n")
	if err := writeSummary(w, httpKeys, r.https); err != nil {
		fmt.Fprintf(os.Stderr, "can not generate HTTP summary: %v", err)
	}
	return nil
}
