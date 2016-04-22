// imirhil.go

/*
  This file contains the datatypes used by tls.imirhil.fr
*/

package main

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

type Source struct {
	Key       Key
	DH        []Key `json:"dh"`
	Protocols []string
	Ciphers   []Cipher
	Score     Score
	HSTS      int       `json:"hsts"`
	Date      time.Time `json:"date"`
}

type Score struct {
	Rank    string
	Details struct {
		Score           float64 `json:"score"`
		Protocol        int     `json:"protocol"`
		KeyExchange     int     `json:"key_exchange"`
		CipherStrengths int     `json:"cipher_strengths"`
	} `json:"details"`
	Error   []interface{}
	Danger  []interface{}
	Warning []interface{}
	Success []string
}

type Report struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	ID      string `json:"_id"`
	Version int    `json:"_version"`
	Found   bool
	Source  Source `json:"_source"`
}
