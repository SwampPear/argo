package settings

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"gopkg.in/yaml.v3"
)

// Supported HTTP methods.
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

// Supported CAPTCHA policies.
type CaptchaPolicy int
const (
	CaptchaManual CaptchaPolicy = iota
)

var captchaPolicyFromStr = map[string]CaptchaPolicy{
	"manual": CaptchaManual,
}

func (m *HTTPMethod) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err == nil {
		if v, ok := httpMethodFromStr[strings.ToLower(s)]; ok {
			*m = v
			return nil
		}

		return fmt.Errorf("invalid HTTPMethod %q", s)
	}
	
	return fmt.Errorf("invalid HTTPMethod value")
}

func (c *CaptchaPolicy) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err == nil {
		if v, ok := captchaPolicyFromStr[strings.ToLower(s)]; ok {
			*c = v
			return nil
		}
	}

	return fmt.Errorf("invalid CaptchaPolicy value")
}

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
	DestructiveActionsForbid bool          `yaml:"destructive_actions_forbidden"`
	AuthBruteforceForbid     bool          `yaml:"auth_bruteforce_forbidden"`
	DoSForbid                bool          `yaml:"dos_forbidden"`
	CaptchaPolicy            CaptchaPolicy `yaml:"captcha_policy"`
	MaxRequestBodyKB         int           `yaml:"max_request_body_kb"`
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
	Mode        string   `yaml:"mode"`
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

// Default settings.
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