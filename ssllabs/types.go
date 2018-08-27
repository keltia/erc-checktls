// types.go

/*
Package ssllabs These are the types used by SSLLabs/Qualys

This is for API v3
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

// LabsCaaRecord describe the DNS CAA record content
type LabsCaaRecord struct {
	Tag   string
	Value string
	Flags int
}

// LabsCaaPolicy is the policy around CAA usage
type LabsCaaPolicy struct {
	PolicyHostname string          `json:"policyHostname"`
	CaaRecords     []LabsCaaRecord `json:"caaRecords"`
}

// LabsCert describes an X.509 certificate
type LabsCert struct {
	ID                     string
	Subject                string
	CommonNames            []string      `json:"commonNames"`
	AltNames               []string      `json:"altNames"`
	NotBefore              int64         `json:"notBefore"`
	NotAfter               int64         `json:"notAfter"`
	IssuerSubject          string        `json:"issuerSubject"`
	SigAlg                 string        `json:"sigAlg"`
	IssuerLabel            string        `json:"issuerLabel"`
	RevocationInfo         int           `json:"revocationInfo"`
	CrlURIs                []string      `json:"crlURIs"`
	OcspURIs               []string      `json:"ocspURIs"`
	RevocationStatus       int           `json:"revocationStatus"`
	CrlRevocationStatus    int           `json:"crlRevocationStatus"`
	OcspRevocationStatus   int           `json:"ocspRevocationStatus"`
	DnsCaa                 bool          `json:"dnsCaa"`
	CaaPolicy              LabsCaaPolicy `json:"caaPolicy"`
	MustStaple             bool          `json:"mustStaple"`
	Sgc                    int
	ValidationType         string `json:"validationType"`
	Issues                 int
	Sct                    bool
	SHA1Hash               string `json:"sha1Hash"`
	PinSHA256              string `json:"pinSha256"`
	KeyAlg                 string `json:"keyAlg"`
	KeySize                int    `json:"keySize"`
	KeyStrength            int    `json:"keyStrength"`
	KeyKnownDebianInsecure bool   `json:"keyKnownDebianInsecure"`
	Raw                    string `json:"raw"`
}

// LabsChainCert describes the chained certificates
type LabsChainCert struct {
	Subject              string
	Label                string
	NotBefore            int64  `json:"notBefore"`
	NotAfter             int64  `json:"notAfter"`
	IssuerSubject        string `json:"issuerSubject"`
	IssuerLabel          string `json:"issuerLabel"`
	SigAlg               string `json:"sigAlg"`
	Issues               int
	KeyAlg               string `json:"sigAlg"`
	KeySize              int    `json:"keySize"`
	KeyStrength          int    `json:"keyStrength"`
	RevocationStatus     int    `json:"revocationStatus"`
	CrlRevocationStatus  int    `json:"crlRevocationStatus"`
	OcspRevocationStatus int    `json:"ocspRevocationStatus"`
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
	Client         LabsSimClient
	ErrorCode      int    `json:"errorCode"`
	ErrorMessage   string `json:"errorMessage"`
	Attempts       int
	CertChainId    string `json:"certChainId"`
	ProtocolID     int    `json:"protocolId"`
	SuiteID        int    `json:"suiteId"`
	SuiteName      string `json:"suiteName"`
	KxType         string `json:"kxType"`
	KxStrength     int    `json:"kxStrength"`
	DhBits         int    `json:"dhBits"`
	DHP            int    `json:"dhP"`
	DHG            int    `json:"dhG"`
	DHYs           int    `json:"dhYs"`
	NamedGroupBits int    `json:"namedGroupBits"`
	NamedGroupId   int    `json:"namedGroupId"`
	NamedGroupName string `json:"namedGroupName"`
	AlertType      int    `json:"alertType"`
	AlertCode      int    `json:"alertCode"`
	KeyAlg         string `json:"keyAlg"`
	KeySize        int    `json:"keySize"`
	SigAlg         string `json:"sigAlg"`
}

// LabsSimDetails are the result of simulation
type LabsSimDetails struct {
	Results []LabsSimulation
}

// LabsSuite describes a single protocol
type LabsSuite struct {
	ID             int `json:"id"`
	Name           string
	CipherStrength int    `json:"cipherStrength"`
	KxType         string `json:"kxType"`
	KxStrength     int    `json:"kxStrength"`
	DhBits         int    `json:"dhBits"`
	DHP            int    `json:"dhP"`
	DHG            int    `json:"dhG"`
	DHYs           int    `json:"dhYs"`
	NamedGroupBits int    `json:"namedGroupBits"`
	NamedGroupId   int    `json:"namedGroupId"`
	NamedGroudName string `json:"namedGroupName"`
	Q              int
}

// LabsSuites is a set of protocols
type LabsSuites struct {
	Protocol   int
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
	HostName   string `json:"hostName"`
	Status     string
	Error      string
	SourceTime int64 `json:"sourceTime"`
}

// LabsHpkpPin is for pinned keys
type LabsHpkpPin struct {
	HashFunction string `json:"hashFunction"`
	Value        string
}

type LabsHpkpDirective struct {
	Name  string
	Value string
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
	Directives        []LabsHpkpDirective
}

// LabsDrownHost describes a potentially Drown-weak site
type LabsDrownHost struct {
	IP      string `json:"ip"`
	Export  bool
	Port    int
	Special bool
	SSLv2   bool `json:"sslv2"`
	Status  string
}

type LabsCertChain struct {
	ID        string
	CertIds   []string        `json:"certIds"`
	Trustpath []LabsTrustPath `json:"trustpath"`
	Issues    int
	NoSni     bool `json:"noSni"`
}

type LabsTrustPath struct {
	CertIds       []string    `json:"certIds"`
	Trust         []LabsTrust `json:"trust"`
	IsPinned      bool        `json:"isPinned"`
	MatchedPins   int         `json:"matchedPins"`
	UnMatchedPins int         `json:"unMatchedPins"`
}

type LabsTrust struct {
	RootStore         string `json:"rootStore"`
	IsTrusted         bool   `json:"isTrusted"`
	TrustErrorMessage string `json:"trustErrorMessage"`
}

type LabsNamedGroups struct {
	List       []LabsNamedGroup
	Preference bool
}

type LabsNamedGroup struct {
	ID   int
	Name string
	Bits int
}

type LabsHttpTransaction struct {
	RequestUrl        string           `json:"requestUrl"`
	StatusCode        int              `json:"statusCode"`
	RequestLine       string           `json:"requestLine"`
	RequestHeaders    []string         `json:"requestHeaders"`
	ResponseLine      string           `json:"responseLine"`
	ResponseRawHeader []string         `json:"responseRawHeader"`
	ResponseHeader    []LabsHttpHeader `json:"responseHeader"`
	FragileServer     bool             `json:"fragileServer"`
}

type LabsHttpHeader struct {
	Name  string
	Value string
}

// LabsEndpointDetails gives the details of a given Endpoint
type LabsEndpointDetails struct {
	HostStartTime                  int64           `json:"hostStartTime"`
	CertChains                     []LabsCertChain `json:"certChains"`
	Protocols                      []LabsProtocol
	Suites                         []LabsSuites
	NoSniSuites                    LabsSuites      `json:"noSniSuites"`
	NamedGroups                    LabsNamedGroups `json:"named_groups"`
	ServerSignature                string          `json:"serverSignature"`
	PrefixDelegation               bool            `json:"prefixDelegation"`
	NonPrefixDelegation            bool            `json:"nonPrefixDelegation"`
	VulnBeast                      bool            `json:"vulnBeast"`
	RenegSupport                   int             `json:"renegSupport"`
	SessionResumption              int             `json:"sessionResumption"`
	CompressionMethods             int             `json:"compressionMethods"`
	SupportsNpn                    bool            `json:"supportsNpn"`
	NpnProcotols                   string          `json:"npnProtocols"`
	SupportsAlpn                   bool            `json:"supportsAlpn"`
	AlpnProtocols                  string
	SessionTickets                 int    `json:"sessionTickets"`
	OcspStapling                   bool   `json:"ocspStapling"`
	StaplingRevocationStatus       int    `json:"staplingRevocationStatus"`
	StaplingRevocationErrorMessage string `json:"staplingRevocationErrorMessage"`
	SniRequired                    bool   `json:"sniRequired"`
	HTTPStatusCode                 int    `json:"httpStatusCode"`
	HTTPForwarding                 string `json:"httpForwarding"`
	SupportsRC4                    bool   `json:"supportsRc4"`
	RC4WithModern                  bool   `json:"rc4WithModern"`
	RC4Only                        bool   `json:"rc4Only"`
	ForwardSecrecy                 int    `json:"forwardSecrecy"`
	SupportAead                    bool   `json:"supportsAead"`
	ProtocolIntolerance            int    `json:"protocolIntolerance"`
	MiscIntolerance                int    `json:"miscIntolerance"`
	Sims                           LabsSimDetails
	Heartbleed                     bool
	Heartbeat                      bool
	OpenSSLCcs                     int `json:"openSslCcs"`
	OpenSSLLuckyMinus20            int `json:"openSSLLuckyMinus20"`
	Ticketbleed                    int
	Bleichenbacher                 int
	Poodle                         bool
	PoodleTLS                      int  `json:"poodleTLS"`
	FallbackScsv                   bool `json:"fallbackScsv"`
	Freak                          bool
	HasSct                         int      `json:"hasSct"`
	DhPrimes                       []string `json:"dhPrimes"`
	DhUsesKnownPrimes              int      `json:"dhUsesKnownPrimes"`
	DhYsReuse                      bool     `json:"dhYsReuse"`
	EcdhParameterReuse             bool     `json:"ecdhParameterReuse"`
	Logjam                         bool
	ChaCha20Preference             bool
	HstsPolicy                     LabsHstsPolicy    `json:"hstsPolicy"`
	HstsPreloads                   []LabsHstsPreload `json:"hstsPreloads"`
	HpkpPolicy                     LabsHpkpPolicy    `json:"hpkpPolicy"`
	HpkpRoPolicy                   LabsHpkpPolicy    `json:"hpkpRoPolicy"`
	DrownHosts                     []interface{}     `json:"drownHosts"`
	DrownErrors                    bool              `json:"drownErrors"`
	DrownVulnerable                bool              `json:"drownVulnerable"`
}

// LabsEndpoint is an Endpoint
type LabsEndpoint struct {
	IPAddress            string `json:"ipAddress"`
	ServerName           string `json:"serverName"`
	StatusMessage        string `json:"statusMessage"`
	StatusDetailsMessage string `json:"statusDetailsMessage"`
	Grade                string
	GradeTrustIgnored    string `json:"gradeTrustIgnored"`
	FutureGrade          string
	HasWarnings          bool `json:"hasWarnings"`
	IsExceptional        bool `json:"isExceptional"`
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
	StatusMessage   string   `json:"statusMessage"`
	StartTime       int64    `json:"startTime"`
	TestTime        int64    `json:"testTime"`
	EngineVersion   string   `json:"engineVersion"`
	CriteriaVersion string   `json:"criteriaVersion"`
	CacheExpiryTime int64    `json:"cacheExpiryTime"`
	CertHostnames   []string `json:"certHostnames"`
	Endpoints       []LabsEndpoint
	Certs           []LabsCert
	RawJSON         string `json:"rawJson"`
}

// LabsReports is a shortcut to all reports
type LabsReports []LabsReport

// LabsResults are all the result of a run w/ 1 or more sites
type LabsResults struct {
	reports   []LabsReport
	responses []string
}

// LabsInfo describes the current SSLLabs engine used
type LabsInfo struct {
	EngineVersion        string `json:"engineVersion"`
	CriteriaVersion      string `json:"criteriaVersion"`
	MaxAssessments       int    `json:"maxAssessments"`
	CurrentAssessments   int    `json:"currentAssessments"`
	NewAssessmentCoolOff int64  `json:"newAssessmentCoolOff"`
	Messages             []string
}
