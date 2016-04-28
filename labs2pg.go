// labs2pg.go

/*
This package implements reading the json from ssllabs-scan output
into our Pg database.
 */
package main

import (
	"flag"

	"github.com/keltia/erc-checktls/ssllabs"
//	"github.com/astaxie/beego/orm"
    _ "github.com/lib/pq" // import your used driver
	"fmt"
)

var (
	contracts map[string]string
)

// init is for pg connection and stuff
func init() {
    // set default database
    //orm.RegisterDataBase("default", "postgres", "roberto", 30)
}

// main is the the starting point
func main() {
	flag.Parse()

	file := flag.Arg(0)

	raw, err := getResults(file)
	if err != nil {
		panic("Can't read " + file)
	}

	// raw is the []byte array to be deserialized into LabsReports
	allSites, err := ssllabs.ParseResults(raw)
	if err != nil {
		panic("Can't parse " + string(raw) + ":" + err.Error())
	}

	// We need that for the reports
	contracts, err = readContractFile("sites-list.csv")

	// generate the final report
	final, err := NewTLSReport(allSites)

	// XXX Early debugging
	fmt.Printf("%v\n", final)

}
