package settings

// llm
type LLMSettings struct {
	Provider            string
	AllowInternetAccess bool
	MaxContextTokens    int
}

// program
type ScopeProgramSettings struct {
	Client  string
	Contact string
	Notes   string
}

// scope
type HTTPMethod int

const (
	HTTPGet HTTPMethod = iota
	HTTPHead
	HTTPOptions
	HTTPPost
	HTTPPut
	HTTPPatch
	HTTPDelete
)

type CaptchaPolicy int

const (
	CaptchaManual CaptchaPolicy = iota
)

type ScopeRulesSettings struct {
	AllowedHTTPMethods        []HTTPMethod
	DestructiveActionsForbid  bool // never run exploits that modify data at scale
	AuthBruteforceForbid      bool // never run expensive brute force operations
	DoSForbid                 bool // includes rate abuse and resource exhaustion
	CaptchaPolicy             CaptchaPolicy
	MaxRequestBodyKB          int // guard against oversized payloads
}

// compliance
type ComplianceLoggingSettings struct {
	ImmutableAuditLog     bool
	IncludeRequestBodies  bool
	IncludeResponseBodies bool
}

type ComplianceNotificationSettings struct {
	OnViolation  []string
	OnCompletion []string
}

type ComplianceSettings struct {
	SafeHarbor    bool
	Logging       ComplianceLoggingSettings
	Notifications ComplianceNotificationSettings
}

// rate limits
type RateLimitGlobal struct {
	RequestsPerMinute  uint
	ConcurrentRequests uint
}

type PerHost struct {
	Host               string
	RequestsPerMinute  uint
	ConcurrentRequests uint
}

type RateLimitSettings struct {
	Global  RateLimitGlobal
	PerHost []PerHost
}

// assets
type Asset struct {
	Mode        string   // e.g., "web", "api", "mobile-api"
	Hostname    string
	Paths       []string
	Ports       []uint
	Schemes     []string // e.g., ["https"]
	Description string
}

type AssetSettings struct {
	InScope     []Asset
	OutOfScope  []Asset
}

// authentication
type AuthenticationAccount struct {
	Username string
	Role     string
	Password string
}

type AuthenticationSettings struct {
	Allowed               bool
	TestAccountsProvided  bool
	Accounts              []AuthenticationAccount
}

type Settings struct {
	LLM            LLMSettings
	Program        ScopeProgramSettings
	Rules          ScopeRulesSettings
	Compliance     ComplianceSettings
	RateLimits     RateLimitSettings
	Assets         AssetSettings
	Authentication AuthenticationSettings
}

// Default returns sensible baseline settings.
func Default() Settings {
	return Settings{
		LLM: LLMSettings{
			Provider:            "local",
			AllowInternetAccess: false,
			MaxContextTokens:    32000,
		},
		Program: ScopeProgramSettings{
			Client:  "",
			Contact: "",
			Notes:   "",
		},
		Rules: ScopeRulesSettings{
			AllowedHTTPMethods: []HTTPMethod{
				HTTPGet, HTTPHead, HTTPOptions, HTTPPost, HTTPPut, HTTPPatch, HTTPDelete,
			},
			DestructiveActionsForbid: false,
			AuthBruteforceForbid:     false,
			DoSForbid:                false,
			CaptchaPolicy:            CaptchaManual,
			MaxRequestBodyKB:         256,
		},
		Compliance: ComplianceSettings{
			SafeHarbor: true,
			Logging: ComplianceLoggingSettings{
				ImmutableAuditLog:     true,
				IncludeRequestBodies:  true,
				IncludeResponseBodies: true,
			},
			Notifications: ComplianceNotificationSettings{},
		},
		RateLimits: RateLimitSettings{
			Global: RateLimitGlobal{
				RequestsPerMinute:  60,
				ConcurrentRequests: 5,
			},
			PerHost: nil,
		},
		Assets: AssetSettings{},
		Authentication: AuthenticationSettings{
			Allowed:              true,
			TestAccountsProvided: true,
			Accounts: nil,
		},
	}
}
