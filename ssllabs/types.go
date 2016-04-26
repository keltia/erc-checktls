// types.go

/*
These are the types used by SSLLabs/Qualys
*/
package ssllabs

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
	DebianFlaw bool `json:"debianFlaw"`
	Q          int
}

// LabsCert describes an X.509 certificate
type LabsCert struct {
	Subject              string
	CommonNames          []string `json:"commonNames"`
	AltNames             []string `json:"altNames"`
	NotBefore            int64    `json:"notBefore"`
	NotAfter             int64    `json:"notAfter"`
	IssuerSubject        string   `json:"issuerSubject"`
	SigAlg               string   `json:"sigAlg"`
	IssuerLabel          string   `json:"issuerLabel"`
	RevocationInfo       int      `json:"revocationInfo"`
	CrlURIs              []string `json:"crlURIs"`
	OcspURIs             []string `json:"ocspUTIs"`
	RevocationStatus     int      `json:"revocationStatus"`
	CrlRevocationStatus  int      `json:"crlRevocationStatus"`
	OcspRevocationStatus int      `json:"ocspRevocationStatus"`
	Sgc                  int
	ValidationType       string `json:"validationType"`
	Issues               int
	Sct                  bool
	SHA1Hash             string `json:"sha1Hash"`
	PinSHA256            string `json:"pinSha256"`
}

// LabsChainCert describes the chained certificates
type LabsChainCert struct {
	Subject              string
	Label                string
	NotBefore            int64  `json:"notBefore"`
	NotAfter             int64  `json:"notAfter"`
	IssuerSubject        string `json:"issuerSubject"`
	SigAlg               string `json:"sigAlg"`
	IssuerLabel          string `json:"issuerLabel"`
	Issues               int
	KeyAlg               string `json:"sigAlg"`
	KeySize              int    `json:"keySize"`
	KeyStrength          int    `json:"keyStrength"`
	RevocationStatus     int    `json:"revocationStatus"`
	CrlRevocationStatus  int    `json:"crlRevocationStatus"`
	OcspRevocationStatus int    `json:"ocspRevocationStatus"`
	SHA1Hash             string `json:"sha1Hash"`
	PinSHA256            string `json:"pinSha256"`
	Raw                  string
}

// LabsChain is the certificate chain envelope
type LabsChain struct {
	Certs  []LabsChainCert
	Issues int
}

// LabsProtocol describes the HTTP protocols
type LabsProtocol struct {
	ID               int `json:"id"`
	Name             string
	Version          string
	V2SuitesDisabled bool
	ErrorMessage     bool `json:"errorMessage"`
	Q                int
}

// LabsSimClient is a simulated client
type LabsSimClient struct {
	ID          int `json:"id"`
	Name        string
	Platform    string
	Version     string
	IsReference bool `json:"isReference"`
}

// LabsSimulation describes the simulation of a given client
type LabsSimulation struct {
	Client     LabsSimClient
	ErrorCode  int `json:"errorCode"`
	Attempts   int
	ProtocolID int `json:"protocolId"`
	SuiteID    int `json:"suiteId"`
}

// LabsSimDetails are the result of simulation
type LabsSimDetails struct {
	Results []LabsSimulation
}

// LabsSuite describes a single protocol
type LabsSuite struct {
	ID             int `json:"id"`
	Name           string
	CipherStrength int `json:"cipherStrength"`
	DHStrength     int `json:"dhStrength"`
	DHP            int `json:"dhP"`
	DHG            int `json:"dhG"`
	DHYs           int `json:"dhYs"`
	ECDHBits       int `json:"ecdhBits"`
	ECDHStrength   int `json:"ecdhStrength"`
	Q              int
}

// LabsSuites is a set of protocols
type LabsSuites struct {
	List       []LabsSuite
	Preference bool
}

func (ls *LabsSuites) len() int {
	return len(ls.List)
}

// LabsHstsPolicy describes the HSTS policy
type LabsHstsPolicy struct {
	LongMaxAge        int64 `json:"LONG_MAX_AGE"`
	Header            string
	Status            string
	Error             string
	MaxAge            int64 `json:"maxAge"`
	IncludeSubDomains bool  `json:"includeSubDomains"`
	Preload           bool
	Directives        map[string]string
}

// LabsHstsPreload is for HSTS preloading
type LabsHstsPreload struct {
	Source     string
	Status     string
	Error      string
	SourceTime int64 `json:"sourceTime"`
}

// LabsHpkpPin is for pinned keys
type LabsHpkpPin struct {
	HashFunction string `json:"hashFunction"`
	Value        string
}

// LabsHpkpPolicy describes the HPKP policy
type LabsHpkpPolicy struct {
	Header            string
	Status            string
	Error             string
	MaxAge            int64 `json:"maxAge"`
	IncludeSubDomains bool  `json:"includeSubDomains"`
	ReportURI         string
	Pins              []LabsHpkpPin
	MatchedPins       []LabsHpkpPin `json:"matchedPins"`
	Directives        map[string]string
}

// LabsEndpointDetails gives the details of a given Endpoint
type LabsEndpointDetails struct {
	HostStartTime            int64 `json:"hostStartTime"`
	Key                      LabsKey
	Cert                     LabsCert
	Chain                    LabsChain
	Protocols                []LabsProtocol
	Suites                   LabsSuites
	ServerSignature          string `json:"serverSignature"`
	PrefixDelegation         bool   `json:"prefixDelegation"`
	NonPrefixDelegation      bool   `json:"nonPrefixDelegation"`
	VulnBeast                bool   `json:"vulnBeast"`
	RenegSupport             int    `json:"renegSupport"`
	StsStatus                string `json:"stsStatus"`
	StsResponseHeader        string `json:"stsResponseHeader"`
	StsMaxAge                int    `json:"stsMaxAge"`
	StsSubdomains            bool   `json:"stsSubdomains"`
	StsPreload               bool   `json:"stsPreload"`
	SessionResumption        int    `json:"sessionResumption"`
	CompressionMethods       int    `json:"compressionMethods"`
	SupportsNpn              bool   `json:"supportsNpn"`
	SessionTickets           int    `json:"sessionTickets"`
	OcspStapling             bool   `json:"ocspStapling"`
	StaplingRevocationStatus int    `json:"staplingRevocationStatus"`
	SniRequired              bool   `json:"sniRequired"`
	HTTPStatusCode           int    `json:"httpStatusCode"`
	SupportsRC4              bool   `json:"supportsRc4"`
	RC4WithModern            bool   `json:"rc4WithModern"`
	RC4Only                  bool   `json:"rc4Only"`
	ForwardSecrecy           int    `json:"forwardSecrecy"`
	Sims                     LabsSimDetails
	Heartbleed               bool
	Heartbeat                bool
	OpenSSLCcs               int `json:"openSslCcs"`
	Poodle                   bool
	PoodleTLS                int  `json:"poodleTLS"`
	FallbackScsv             bool `json:"fallbackScsv"`
	Freak                    bool
	HasSct                   int      `json:"hasSct"`
	DhPrimes                 []string `json:"dhPrimes"`
	DhUsesKnownPrimes        int      `json:"dhUsesKnownPrimes"`
	DhYsReuse                bool     `json:"dhYsReuse"`
	Logjam                   bool
	ChaCha20Preference       bool
	HstsPolicy               LabsHstsPolicy    `json:"hstsPolicy"`
	HstsPreloads             []LabsHstsPreload `json:"hstsPreloads"`
	HpkpPolicy               LabsHpkpPolicy    `json:"hpkpPolicy"`
	HpkpRoPolicy             LabsHpkpPolicy    `json:"hpkpRoPolicy"`
}

// LabsEndpoint is an Endpoint
type LabsEndpoint struct {
	IPAddress            string `json:"ipAddress"`
	ServerName           string `json:"serverName"`
	StatusMessage        string `json:"statusMessage"`
	StatusDetailsMessage string `json:"statusDetailsMessage"`
	Grade                string
	GradeTrustIgnored    string `json:"gradeTrustIgnored"`
	HasWarnings          bool   `json:"hasWarnings"`
	IsExceptional        bool   `json:"isExceptional"`
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
	IsPublic        bool `json:"isPublic"`
	Status          string
	StatusMessage   string `json:"statusMessage"`
	StartTime       int64  `json:"startTime"`
	TestTime        int64  `json:"testTime"`
	EngineVersion   string `json:"engineVersion"`
	CriteriaVersion string `json:"criteriaVerison"`
	CacheExpiryTime int64  `json:"cacheExpiryTime"`
	Endpoints       []LabsEndpoint
	CertHostnames   []string `json:"certHostnames"`
	rawJSON         string   `json:"rawJson"`
}

// LabsResults are all the result of a run w/ 1 or more sites
type LabsResults struct {
	reports   []LabsReport
	responses []string
}

// LabsReports is a shortcut to all reports
type LabsReports []LabsReport

// LabsInfo describes the current SSLLabs engine used
type LabsInfo struct {
	EngineVersion        string `json:"engineVersion"`
	CriteriaVersion      string `json:"criteriaVersion"`
	MaxAssessments       int    `json:"maxAssessments"`
	CurrentAssessments   int    `json:"currentAssessments"`
	NewAssessmentCoolOff int64  `json:"newAssessmentCoolOff"`
	Messages             []string
}
