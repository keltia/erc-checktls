// labs2pg.go

/*
This package implements reading the json from ssllabs-scan output
into our Pg database.
 */
package main

import (
	"flag"

//	"github.com/astaxie/beego/orm"
    _ "github.com/lib/pq" // import your used driver
)

var (
	contracts map[string]string
	reports   *[]LabsReport
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

	reports, err := parseResults(raw)
	if err != nil {
		panic("Can't parse " + string(raw) + ":" + err.Error())
	}

	contracts, err = readContractFile("sites-list.csv")
	err = insertResults(reports)
}
