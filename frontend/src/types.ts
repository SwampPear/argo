// Enums (JSON uses strings)
export type HTTPMethod = 'GET' | 'HEAD' | 'OPTIONS' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
export type CaptchaPolicy = 'manual'

// LLM
export interface LLMSettings {
  provider: string
  allow_internet_access: boolean
  max_context_tokens: number
}

// Program
export interface ScopeProgramSettings {
  client: string
  contact: string
  notes: string
}

// Scope rules
export interface ScopeRulesSettings {
  allowed_http_methods: HTTPMethod[]
  destructive_actions_forbidden: boolean
  auth_bruteforce_forbidden: boolean
  dos_forbidden: boolean
  captcha_policy: CaptchaPolicy
  max_request_body_kb: number
}

// Compliance
export interface ComplianceLoggingSettings {
  immutable_audit_log: boolean
  include_request_bodies: boolean
  include_response_bodies: boolean
}

export interface ComplianceNotificationSettings {
  on_violation: string[]
  on_completion: string[]
}

export interface ComplianceSettings {
  safe_harbor: boolean
  logging: ComplianceLoggingSettings
  notifications: ComplianceNotificationSettings
}

// Rate limits
export interface RateLimitGlobal {
  requests_per_minute: number
  concurrent_requests: number
}

export interface PerHost {
  host: string
  requests_per_minute: number
  concurrent_requests: number
}

export interface RateLimitSettings {
  global: RateLimitGlobal
  per_host: PerHost[]
}

// Assets
export interface Asset {
  mode: string        // "web" | "api" | "mobile-api" (free-form here)
  hostname: string
  paths: string[]
  ports: number[]
  schemes: string[]   // e.g., ["https"]
  description: string
}

export interface AssetSettings {
  in_scope: Asset[]
  out_of_scope: Asset[]
}

// Authentication
export interface AuthenticationAccount {
  username: string
  role: string
  password: string
}

export interface AuthenticationSettings {
  allowed: boolean
  test_accounts_provided: boolean
  accounts: AuthenticationAccount[]
}

// Root
export interface Settings {
  llm: LLMSettings
  program: ScopeProgramSettings
  rules: ScopeRulesSettings
  compliance: ComplianceSettings
  rate_limits: RateLimitSettings
  assets: AssetSettings
  authentication: AuthenticationSettings
}
