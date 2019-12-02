// types.go

package site

const (
	// 1, 2, 3 are the main categories 1=green, 2=yellow, 3=red
	CatGreen = 1 + iota
	CatYellow
	CatRed

	// 1 is for correct https w/ redirection, 2 is https&http, 3 is http only
	TypeError = 1 + iota
	TypeHTTPSok
	TypeHTTPSnok
	TypeHTTP
)

// TLSSite is a summary for each site
type TLSSite struct {
	Name     string
	Contract string

	Empty bool

	Grade      string
	CryptCheck string
	Mozilla    string

	DefKey bool
	DefSig bool
	DefCA  string

	IsExpired bool
	Issues    bool

	Protocols string
	PFS       bool

	OCSPStapling bool
	HSTS         int64
	Sweet32      bool

	Type    int
	CatHTTP int
	CatTLS  int

	// Could we connect at all?
	Connect bool
}

type Flags struct {
	IgnoreMozilla bool
	IgnoreImirhil bool
	LogLevel      int
	Contracts     map[string]string
}
