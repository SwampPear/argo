package settings

import (
	"bytes"
	"fmt"
	"os"
	"io"

	"gopkg.in/yaml.v3"
)

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

var httpMethodFromStr = map[string]HTTPMethod{
	"get":     HTTPGet,
	"head":    HTTPHead,
	"options": HTTPOptions,
	"post":    HTTPPost,
	"put":     HTTPPut,
	"patch":   HTTPPatch,
	"delete":  HTTPDelete,
}
var httpMethodToStr = []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}

func (m HTTPMethod) String() string {
	if int(m) >= 0 && int(m) < len(httpMethodToStr) {
		return httpMethodToStr[m]
	}
	return "UNKNOWN"
}

func (m *HTTPMethod) UnmarshalYAML(node *yaml.Node) error {
	// accept either string (case-insensitive) or integer
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("HTTPMethod must be a scalar")
	}
	var asStr string
	if err := node.Decode(&asStr); err == nil {
		if v, ok := httpMethodFromStr[lower(asStr)]; ok {
			*m = v
			return nil
		}
		return fmt.Errorf("invalid HTTPMethod %q", asStr)
	}
	var asInt int
	if err := node.Decode(&asInt); err == nil {
		*m = HTTPMethod(asInt)
		return nil
	}
	return fmt.Errorf("invalid HTTPMethod value")
}

type CaptchaPolicy int

const (
	CaptchaManual CaptchaPolicy = iota
)

func (c *CaptchaPolicy) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("CaptchaPolicy must be a scalar")
	}
	var s string
	if err := node.Decode(&s); err == nil {
		switch lower(s) {
		case "manual", "captcha_manual", "0":
			*c = CaptchaManual
			return nil
		default:
			return fmt.Errorf("invalid CaptchaPolicy %q", s)
		}
	}
	var i int
	if err := node.Decode(&i); err == nil {
		*c = CaptchaPolicy(i)
		return nil
	}
	return fmt.Errorf("invalid CaptchaPolicy value")
}

func lower(s string) string {
	// tiny helper to avoid importing strings for one call
	b := []byte(s)
	for i := range b {
		if 'A' <= b[i] && b[i] <= 'Z' {
			b[i] = b[i] + 32
		}
	}
	return string(b)
}

// ====== Structs with YAML tags ======

type LLMSettings struct {
	Provider            string `yaml:"provider"`
	AllowInternetAccess bool   `yaml:"allow_internet_access"`
	MaxContextTokens    int    `yaml:"max_context_tokens"`
}

type ScopeProgramSettings struct {
	Client  string `yaml:"client"`
	Contact string `yaml:"contact"`
	Notes   string `yaml:"notes"`
}

type ScopeRulesSettings struct {
	AllowedHTTPMethods       []HTTPMethod  `yaml:"allowed_http_methods"`
	DestructiveActionsForbid bool         `yaml:"destructive_actions_forbidden"`
	AuthBruteforceForbid     bool         `yaml:"auth_bruteforce_forbidden"`
	DoSForbid                bool         `yaml:"dos_forbidden"`
	CaptchaPolicy            CaptchaPolicy `yaml:"captcha_policy"`
	MaxRequestBodyKB         int          `yaml:"max_request_body_kb"`
}

type ComplianceLoggingSettings struct {
	ImmutableAuditLog     bool `yaml:"immutable_audit_log"`
	IncludeRequestBodies  bool `yaml:"include_request_bodies"`
	IncludeResponseBodies bool `yaml:"include_response_bodies"`
}

type ComplianceNotificationSettings struct {
	OnViolation  []string `yaml:"on_violation"`
	OnCompletion []string `yaml:"on_completion"`
}

type ComplianceSettings struct {
	SafeHarbor    bool                           `yaml:"safe_harbor"`
	Logging       ComplianceLoggingSettings      `yaml:"logging"`
	Notifications ComplianceNotificationSettings `yaml:"notifications"`
}

type RateLimitGlobal struct {
	RequestsPerMinute  uint `yaml:"requests_per_minute"`
	ConcurrentRequests uint `yaml:"concurrent_requests"`
}

type PerHost struct {
	Host               string `yaml:"host"`
	RequestsPerMinute  uint   `yaml:"requests_per_minute"`
	ConcurrentRequests uint   `yaml:"concurrent_requests"`
}

type RateLimitSettings struct {
	Global  RateLimitGlobal `yaml:"global"`
	PerHost []PerHost       `yaml:"per_host"`
}

type Asset struct {
	Mode        string   `yaml:"mode"` // e.g., "web", "api", "mobile-api"
	Hostname    string   `yaml:"hostname"`
	Paths       []string `yaml:"paths"`
	Ports       []uint   `yaml:"ports"`
	Schemes     []string `yaml:"schemes"`
	Description string   `yaml:"description"`
}

type AssetSettings struct {
	InScope    []Asset `yaml:"in_scope"`
	OutOfScope []Asset `yaml:"out_of_scope"`
}

type AuthenticationAccount struct {
	Username string `yaml:"username"`
	Role     string `yaml:"role"`
	Password string `yaml:"password"`
}

type AuthenticationSettings struct {
	Allowed              bool                   `yaml:"allowed"`
	TestAccountsProvided bool                   `yaml:"test_accounts_provided"`
	Accounts             []AuthenticationAccount `yaml:"accounts"`
}

type Settings struct {
	LLM            LLMSettings            `yaml:"llm"`
	Program        ScopeProgramSettings   `yaml:"program"`
	Rules          ScopeRulesSettings     `yaml:"rules"`
	Compliance     ComplianceSettings     `yaml:"compliance"`
	RateLimits     RateLimitSettings      `yaml:"rate_limits"`
	Assets         AssetSettings          `yaml:"assets"`
	Authentication AuthenticationSettings `yaml:"authentication"`
}

// Default as you provided (unchanged)
func Default() Settings {
	return Settings{
		LLM: LLMSettings{
			Provider:            "local",
			AllowInternetAccess: false,
			MaxContextTokens:    32000,
		},
		Program: ScopeProgramSettings{},
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
		},
		Assets:         AssetSettings{},
		Authentication: AuthenticationSettings{Allowed: true, TestAccountsProvided: true},
	}
}


// LoadYAML overlays file values onto Default() so missing keys keep defaults.
func LoadYAML(path string) (Settings, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read %s: %w", path, err)
	}

	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true) // strict: error on unknown keys
	if err := dec.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("parse %s: %w", path, err)
	}
	return cfg, nil
}

// ---- tiny helper to use yaml.Decoder with []byte without importing bytes in user code

type byteReader struct{ b []byte; i int }

func bytesReader(b []byte) *yaml.Decoder {
	// emulate io.Reader for yaml.Decoder
	return yaml.NewDecoder(&byteReader{b: b})
}
func (r *byteReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
