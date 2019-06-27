package main

import (
	"fmt"
	"strings"

	"github.com/keltia/ssllabs"
)

func displayWildcards(all []ssllabs.Host) string {
	var buf strings.Builder

	fmt.Fprint(&buf, "")
	// Now analyze each site
	for _, site := range all {
		debug("site=%s\n", site.Host)
		if site.Endpoints != nil {
			// If Certs is empty, we could not connect
			if len(site.Certs) != 0 {
				cert := site.Certs[0]
				debug("  cert=%#v\n", cert)
				debug("  CN=%s\n", cert.Subject)

				if strings.HasPrefix(cert.Subject, "CN=*") {
					debug("adding %s\n", site.Host)
					buf.WriteString(fmt.Sprintf("  %-35s %-16s CN=%v \n", site.Host, site.Endpoints[0].IPAddress, cert.CommonNames))
				}
			}
		}
	}
	return buf.String()
}

