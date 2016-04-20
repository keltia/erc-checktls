// types.go

package main

import (
	"encoding/json"
)

// LabsError is for whatever error we get from SSLLabs
type LabsError struct {
	Field   string
	Message string
}

// LabsErrorResponse is a set of errors
type LabsErrorResponse struct {
	ResponseErrors []LabsError `json:"errors"`
}

// Error() implements the interface
func (e LabsErrorResponse) Error() string {
	msg, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(msg)
}

// LabsKey describes a crypto key
type LabsKey struct {
	Size       int
	Strength   int
	Alg        string
	DebianFlaw bool
	Q          int
}

// LabsCert describes an X.509 certificate
type LabsCert struct {
	Subject              string
	CommonNames          []string
	AltNames             []string
	NotBefore            int64
	NotAfter             int64
	IssuerSubject        string
	SigAlg               string
	IssuerLabel          string
	RevocationInfo       int
	CrlURIs              []string
	OcspURIs             []string
	RevocationStatus     int
	CrlRevocationStatus  int
	OcspRevocationStatus int
	Sgc                  int
	ValidationType       string
	Issues               int
	Sct                  bool
}

// LabsChainCert describes the chained certificates
type LabsChainCert struct {
	Subject              string
	Label                string
	NotBefore            int64
	NotAfter             int64
	IssuerSubject        string
	IssuerLabel          string
	SigAlg               string
	Issues               int
	KeyAlg               string
	KeySize              int
	KeyStrength          int
	RevocationStatus     int
	CrlRevocationStatus  int
	OcspRevocationStatus int
	Raw                  string
}

// LabsChain is the certificate chain envelope
type LabsChain struct {
	Certs  []LabsChainCert
	Issues int
}

// LabsProtocol describes the HTTP protocols
type LabsProtocol struct {
	ID               int
	Name             string
	Version          string
	V2SuitesDisabled bool
	ErrorMessage     bool
	Q                int
}

// LabsSimClient is a simulated client
type LabsSimClient struct {
	ID          int
	Name        string
	Platform    string
	Version     string
	IsReference bool
}

// LabsSimulation describes the simulation of a given client
type LabsSimulation struct {
	Client     LabsSimClient
	ErrorCode  int
	Attempts   int
	ProtocolID int
	SuiteID    int
}

// LabsSimDetails are the result of simulation
type LabsSimDetails struct {
	Results []LabsSimulation
}

// LabsSuite describes a single protocol
type LabsSuite struct {
	ID             int
	Name           string
	CipherStrength int
	DhStrength     int
	DhP            int
	DhG            int
	DhYs           int
	EcdhBits       int
	EcdhStrength   int
	Q              int
}

// LabsSuites is a set of protocols
type LabsSuites struct {
	List       []LabsSuite
	Preference bool
}

// LabsHstsPolicy describes the HSTS policy
type LabsHstsPolicy struct {
	LONG_MAX_AGE      int64
	Header            string
	Status            string
	Error             string
	MaxAge            int64
	IncludeSubDomains bool
	Preload           bool
	Directives        map[string]string
}

// LabsHstsPreload is for HSTS preloading
type LabsHstsPreload struct {
	Source     string
	Status     string
	Error      string
	SourceTime int64
}

// LabsHpkpPin is for pinned keys
type LabsHpkpPin struct {
	HashFunction string
	Value        string
}

// LabsHpkpPolicy describes the HPKP policy
type LabsHpkpPolicy struct {
	Header            string
	Status            string
	Error             string
	MaxAge            int64
	IncludeSubDomains bool
	ReportURI         string
	Pins              []LabsHpkpPin
	MatchedPins       []LabsHpkpPin
	Directives        map[string]string
}

// LabsEndpointDetails gives the details of a given Endpoint
type LabsEndpointDetails struct {
	HostStartTime                  int64
	Key                            LabsKey
	Cert                           LabsCert
	Chain                          LabsChain
	Protocols                      []LabsProtocol
	Suites                         LabsSuites
	ServerSignature                string
	PrefixDelegation               bool
	NonPrefixDelegation            bool
	VulnBeast                      bool
	RenegSupport                   int
	SessionResumption              int
	CompressionMethods             int
	SupportsNpn                    bool
	NpnProtocols                   string
	SessionTickets                 int
	OcspStapling                   bool
	StaplingRevocationStatus       int
	StaplingRevocationErrorMessage string
	SniRequired                    bool
	HTTPStatusCode                 int
	HTTPForwarding                 string
	ForwardSecrecy                 int
	SupportsRc4                    bool
	Rc4WithModern                  bool
	Rc4Only                        bool
	Sims                           LabsSimDetails
	Heartbleed                     bool
	Heartbeat                      bool
	OpenSslCcs                     int
	Poodle                         bool
	PoodleTLS                      int
	FallbackScsv                   bool
	Freak                          bool
	HasSct                         int
	DhPrimes                       []string
	DhUsesKnownPrimes              int
	DhYsReuse                      bool
	Logjam                         bool
	ChaCha20Preference             bool
	HstsPolicy                     LabsHstsPolicy
	HstsPreloads                   []LabsHstsPreload
	HpkpPolicy                     LabsHpkpPolicy
	HpkpRoPolicy                   LabsHpkpPolicy
}

// LabsEndpoint is an Endpoint
type LabsEndpoint struct {
	IPAddress            string
	ServerName           string
	StatusMessage        string
	StatusDetailsMessage string
	Grade                string
	GradeTrustIgnored    string
	HasWarnings          bool
	IsExceptional        bool
	Progress             int
	Duration             int
	Eta                  int
	Delegation           int
	Details              LabsEndpointDetails
}

// LabsReport is a one-site report
type LabsReport struct {
	Host            string
	Port            int
	Protocol        string
	IsPublic        bool
	Status          string
	StatusMessage   string
	StartTime       int64
	TestTime        int64
	EngineVersion   string
	CriteriaVersion string
	CacheExpiryTime int64
	Endpoints       []LabsEndpoint
	CertHostnames   []string
	rawJSON         string
}

// LabsResults are all the result of a run w/ 1 or more sites
type LabsResults struct {
	reports   []LabsReport
	responses []string
}

// LabsInfo describes the current SSLLabs engine used
type LabsInfo struct {
	EngineVersion        string
	CriteriaVersion      string
	MaxAssessments       int
	CurrentAssessments   int
	NewAssessmentCoolOff int64
	Messages             []string
}
