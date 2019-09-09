package site

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"
)

const (
	DefaultKeySize = 2048
	DefaultAlg     = "RSA"
	DefaultECAlg   = "EC"
	DefaultSig     = "SHA256withRSA"
)

// Use an interface to enable better tests
type Mapi interface {
	GetGrade(string) (string, error)
	IsHTTPSonly(string) (bool, error)
}

type Capi interface {
	GetScore(string) (string, error)
}

type Sapi interface {
	GetDetailedReport(string, ...map[string]string) (ssllabs.Host, error)
}

type Client struct {
	sslc Sapi
	moz  Mapi
	irml Capi

	fnImirhil func(site ssllabs.Host) string
	fnMozilla func(site ssllabs.Host) string
}

var (
	contracts map[string]string

	DefaultIssuer = regexp.MustCompile(`(?i:GlobalSign)`)

	fIgnoreMozilla bool
	fIgnoreImirhil bool
	fVerbose       bool
	fDebug         bool
)

func NewClient(f Flags) *Client {
	c := &Client{}
	contracts = f.Contracts

	fDebug = f.LogLevel >= 2
	fVerbose = f.LogLevel >= 1

	fIgnoreImirhil = f.IgnoreImirhil
	fIgnoreMozilla = f.IgnoreMozilla

	if !fIgnoreImirhil {
		cnf := cryptcheck.Config{
			Log:     f.LogLevel,
			Refresh: true,
			Timeout: 30,
		}
		c.irml = cryptcheck.NewClient(cnf)

		c.fnImirhil = func(site ssllabs.Host) string {
			debug("  imirhil\n")
			score, err := c.irml.GetScore(site.Host)
			if err != nil {
				verbose("cryptcheck error: %s (%s)\n", site.Host, err.Error())
			}
			return score
		}
	} else {
		c.fnImirhil = func(site ssllabs.Host) string {
			return ""
		}
	}

	if !fIgnoreMozilla {
		cnf := observatory.Config{
			Log:     f.LogLevel,
			Timeout: 30,
		}
		c.moz, _ = observatory.NewClient(cnf)

		c.fnMozilla = func(site ssllabs.Host) string {
			debug("  observatory\n")
			score, err := c.moz.GetGrade(site.Host)
			if err != nil {
				verbose("Mozilla error: %s (%s)\n", site.Host, err.Error())
				return "X"
			}
			return score
		}
	} else {
		c.fnMozilla = func(site ssllabs.Host) string {
			return ""
		}
	}

	c.sslc, _ = ssllabs.NewClient()
	return c
}

func (c *Client) New(site string) (TLSSite, error) {
	if site == "" {
		return TLSSite{}, errors.New("Empty site")
	}

	r, err := c.sslc.GetDetailedReport(site)
	return c.NewFromHost(r), errors.Wrap(err, "New")
}

func (c *Client) NewFromHost(site ssllabs.Host) TLSSite {
	var current TLSSite

	debug("NewFromHost")
	if site.Endpoints == nil || len(site.Endpoints) == 0 {
		verbose("Site %s has no endpoint\n", site.Host)
		current = TLSSite{
			Name:     site.Host,
			Contract: contracts[site.Host],
			Empty:    true,
		}
	} else {
		endp := site.Endpoints[0]
		det := endp.Details

		verbose("Host: %s\n", site.Host)

		protos := []string{}
		for _, p := range det.Protocols {
			protos = append(protos, fmt.Sprintf("%sv%s", p.Name, p.Version))
		}

		// FIll in all details
		current = TLSSite{
			Name:         site.Host,
			Contract:     contracts[site.Host],
			Grade:        endp.Grade,
			CryptCheck:   c.fnImirhil(site),
			Mozilla:      c.fnMozilla(site),
			Protocols:    strings.Join(protos, ","),
			PFS:          det.ForwardSecrecy >= 2,
			OCSPStapling: det.OcspStapling,
			HSTS:         checkHSTS(det),
			Sweet32:      checkSweet32(det),
			Type:         c.findServerType(site),
		}

		// Handle case where we have a DNS entry but no connection
		if len(site.Certs) != 0 {
			cert := site.Certs[0]
			current.DefKey = checkKey(cert)

			current.DefCA = checkIssuer(cert, DefaultIssuer)
			current.DefSig = cert.SigAlg == DefaultSig
			current.IsExpired = hasExpired(cert.NotAfter)
		}

		if len(det.CertChains) != 0 {
			current.Issues = det.CertChains[0].Issues != 0
		}
	}
	return current
}

func (c *Client) findServerType(site ssllabs.Host) int {
	// Should be obvious, 2nd field is only present if no valid cert is found
	if len(site.Certs) == 0 || len(site.CertHostnames) != 0 {
		return TypeHTTP
	}

	if !fIgnoreMozilla {
		// Check the Mozilla report
		if yes, _ := c.moz.IsHTTPSonly(site.Host); yes {
			return TypeHTTPSok
		}
	}
	return TypeHTTPSnok
}
