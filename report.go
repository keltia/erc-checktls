// report.go

/*
This file contains func for generating the report
*/
package TLS

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/atotto/encoding/csv"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"

	"github.com/keltia/erc-checktls/site"
)

var (
	logLevel int

	fJobs          int
	fIgnoreImirhil bool
	fIgnoreMozilla bool

	// this is to protect the Sites array
	lock sync.Mutex

	contracts map[string]string
	tmpls     map[string]string
)

func Init(c Config) {
	var err error

	logLevel = c.LogLevel
	fJobs = c.Jobs
	fIgnoreImirhil = c.IgnoreImirhil
	fIgnoreMozilla = c.IgnoreMozilla

	contracts, tmpls, err = loadResources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't load resources %s: %v\n", resourcesPath, err)
	}
}

func getSSLablsVersion(site ssllabs.Host) string {
	debug("%#v", site)
	return fmt.Sprintf("%s/%s", site.EngineVersion, site.CriteriaVersion)
}

func worker(i int, queue chan ssllabs.Host, r *Report, wg *sync.WaitGroup) {
	defer wg.Done()

	// Setup our env.
	c := site.NewClient(site.Flags{
		LogLevel:      logLevel,
		IgnoreMozilla: fIgnoreMozilla,
		IgnoreImirhil: fIgnoreImirhil,
		Contracts:     contracts,
	})

	for payload := range queue {
		verbose("worker%d/starting %s\n", i, payload.Host)

		s := c.NewFromHost(payload)

		verbose("worker%d/inserting %s\n", i, payload.Host)

		// Block on mutex
		lock.Lock()
		r.GatherStats(s)

		r.Sites = append(r.Sites, s)
		lock.Unlock()

		verbose("worker%d/done %s\n", i, payload.Host)
	}
	verbose("worker%d finished\n", i)
}

// NewReport generates everything we need for display/export
func NewReport(reports []ssllabs.Host, jobs int) (r *Report, err error) {

	if len(reports) == 0 {
		return nil, fmt.Errorf("empty list")
	}

	r = &Report{
		Date:    time.Now(),
		SSLLabs: getSSLablsVersion(reports[0]),
		cntrs:   map[string]int{},
		https:   map[string]int{},
	}

	verbose("%d sites found, %d workers.\n", len(reports), jobs)

	queue := make(chan ssllabs.Host, len(reports))

	wg := &sync.WaitGroup{}

	// Setup workers
	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go worker(i, queue, r, wg)

	}

	// Now analyze each site
	for _, s := range reports {
		verbose("queueing %s\n", s.Host)

		queue <- s
	}

	close(queue)
	wg.Wait()

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

func (r *Report) ColourMap(criteria string) Types {
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
func (r *Report) ToCSV(w io.Writer) (err error) {
	wh := csv.NewWriter(w)
	debug("%v\n", r.Sites)
	if err = wh.WriteStructHeader(r.Sites[0]); err != nil {
		return errors.Wrap(err, "can not write csv header")
	}

	err = wh.WriteStructAll(r.Sites)
	return errors.Wrap(err, "can not write csv file")
}

func (r *Report) WriteCSV(w io.Writer) error {
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

func (r *Report) GatherStats(s site.TLSSite) {
	if !s.Empty {
		if s.Grade != "" && s.Grade != "Z" {
			r.cntrs["Total"]++
		} else {
			r.cntrs["Z"]++
		}
		r.cntrs[s.Grade]++
		if s.PFS {
			r.cntrs["PFS"]++
		}
		if s.Sweet32 {
			r.cntrs["Sweet32"]++
		}
		if s.Issues {
			r.cntrs["Issues"]++
		}
		if s.OCSPStapling {
			r.cntrs["OCSPStapling"]++
		}
		if s.HSTS > 0 {
			r.cntrs["HSTS"]++
		}

		if s.Mozilla != "" {
			if s.Mozilla >= "G" {
				r.https["Bad"]++
			} else {
				r.https["Total"]++
			}
			r.https[s.Mozilla]++
		} else {
			r.https["Broken"]++
		}
	} else {
		r.cntrs["X"]++
	}
}
