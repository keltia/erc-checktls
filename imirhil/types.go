// types.go

// XXX Versioning of the API is nonexistent, we have to cope
// 20160510 "old" API
// 20160511 "new" API

package imirhil

import "time"

// Key describes a single key
type Key struct {
	Type    string `json:"type"`
	Size    int    `json:"size"`
	RSASize int    `json:"rsa_size"`
}

// Cipher describes a single cipher
type Cipher struct {
	Protocol string
	Name     string
	Size     int
	DH       Key `json:"dh"`
}

// Grade aka score of the site
type Grade struct {
	Rank    string
	Details struct {
		Score           float64 `json:"score"`
		Protocol        int     `json:"protocol"`
		KeyExchange     int     `json:"key_exchange"`
		CipherStrengths int     `json:"cipher_strengths"`
	} `json:"details"`
	Error   []string
	Danger  []string
	Warning []string
	Success []string
}

// Site contains DNS site data
type Site struct {
	Name string
	IP   string `json:"ip"`
	Port int
}

// Handshake contains crypto parameters
type Handshake struct {
	Key       Key
	DH        []Key `json:"dh"`
	Protocols []string
	Ciphers   []Cipher
	HSTS      int `json:"hsts"`
}

// Host describe a single host
type Host struct {
	Host      Site      `json:"host"`
	Handshake Handshake `json:"handshake"`
	Grade     Grade
}

// Source describes the details for the crypto
type Report struct {
	Hosts []Host
	Date  time.Time `json:"date"`
}
