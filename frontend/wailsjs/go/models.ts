export namespace main {
	
	export class LogEntry {
	    step: number;
	    id: string;
	    timestamp: string;
	    module: string;
	    action: string;
	    target: string;
	    status: string;
	    duration: string;
	    confidence: number;
	    summary: string;
	    parent_step_id: number;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.step = source["step"];
	        this.id = source["id"];
	        this.timestamp = source["timestamp"];
	        this.module = source["module"];
	        this.action = source["action"];
	        this.target = source["target"];
	        this.status = source["status"];
	        this.duration = source["duration"];
	        this.confidence = source["confidence"];
	        this.summary = source["summary"];
	        this.parent_step_id = source["parent_step_id"];
	    }
	}
	export class AppState {
	    projectDir: string;
	    settings: settings.Settings;
	    logs: LogEntry[];
	    version: number;
	
	    static createFrom(source: any = {}) {
	        return new AppState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectDir = source["projectDir"];
	        this.settings = this.convertValues(source["settings"], settings.Settings);
	        this.logs = this.convertValues(source["logs"], LogEntry);
	        this.version = source["version"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace settings {
	
	export class Asset {
	    Mode: string;
	    Hostname: string;
	    Paths: string[];
	    Ports: number[];
	    Schemes: string[];
	    Description: string;
	
	    static createFrom(source: any = {}) {
	        return new Asset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Mode = source["Mode"];
	        this.Hostname = source["Hostname"];
	        this.Paths = source["Paths"];
	        this.Ports = source["Ports"];
	        this.Schemes = source["Schemes"];
	        this.Description = source["Description"];
	    }
	}
	export class AssetSettings {
	    InScope: Asset[];
	    OutOfScope: Asset[];
	
	    static createFrom(source: any = {}) {
	        return new AssetSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.InScope = this.convertValues(source["InScope"], Asset);
	        this.OutOfScope = this.convertValues(source["OutOfScope"], Asset);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AuthenticationAccount {
	    Username: string;
	    Role: string;
	    Password: string;
	
	    static createFrom(source: any = {}) {
	        return new AuthenticationAccount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Username = source["Username"];
	        this.Role = source["Role"];
	        this.Password = source["Password"];
	    }
	}
	export class AuthenticationSettings {
	    Allowed: boolean;
	    TestAccountsProvided: boolean;
	    Accounts: AuthenticationAccount[];
	
	    static createFrom(source: any = {}) {
	        return new AuthenticationSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Allowed = source["Allowed"];
	        this.TestAccountsProvided = source["TestAccountsProvided"];
	        this.Accounts = this.convertValues(source["Accounts"], AuthenticationAccount);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ComplianceLoggingSettings {
	    ImmutableAuditLog: boolean;
	    IncludeRequestBodies: boolean;
	    IncludeResponseBodies: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ComplianceLoggingSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ImmutableAuditLog = source["ImmutableAuditLog"];
	        this.IncludeRequestBodies = source["IncludeRequestBodies"];
	        this.IncludeResponseBodies = source["IncludeResponseBodies"];
	    }
	}
	export class ComplianceNotificationSettings {
	    OnViolation: string[];
	    OnCompletion: string[];
	
	    static createFrom(source: any = {}) {
	        return new ComplianceNotificationSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OnViolation = source["OnViolation"];
	        this.OnCompletion = source["OnCompletion"];
	    }
	}
	export class ComplianceSettings {
	    SafeHarbor: boolean;
	    Logging: ComplianceLoggingSettings;
	    Notifications: ComplianceNotificationSettings;
	
	    static createFrom(source: any = {}) {
	        return new ComplianceSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SafeHarbor = source["SafeHarbor"];
	        this.Logging = this.convertValues(source["Logging"], ComplianceLoggingSettings);
	        this.Notifications = this.convertValues(source["Notifications"], ComplianceNotificationSettings);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LLMSettings {
	    Provider: string;
	    AllowInternetAccess: boolean;
	    MaxContextTokens: number;
	
	    static createFrom(source: any = {}) {
	        return new LLMSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Provider = source["Provider"];
	        this.AllowInternetAccess = source["AllowInternetAccess"];
	        this.MaxContextTokens = source["MaxContextTokens"];
	    }
	}
	export class PerHost {
	    Host: string;
	    RequestsPerMinute: number;
	    ConcurrentRequests: number;
	
	    static createFrom(source: any = {}) {
	        return new PerHost(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Host = source["Host"];
	        this.RequestsPerMinute = source["RequestsPerMinute"];
	        this.ConcurrentRequests = source["ConcurrentRequests"];
	    }
	}
	export class RateLimitGlobal {
	    RequestsPerMinute: number;
	    ConcurrentRequests: number;
	
	    static createFrom(source: any = {}) {
	        return new RateLimitGlobal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RequestsPerMinute = source["RequestsPerMinute"];
	        this.ConcurrentRequests = source["ConcurrentRequests"];
	    }
	}
	export class RateLimitSettings {
	    Global: RateLimitGlobal;
	    PerHost: PerHost[];
	
	    static createFrom(source: any = {}) {
	        return new RateLimitSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Global = this.convertValues(source["Global"], RateLimitGlobal);
	        this.PerHost = this.convertValues(source["PerHost"], PerHost);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ScopeProgramSettings {
	    Client: string;
	    Contact: string;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new ScopeProgramSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Client = source["Client"];
	        this.Contact = source["Contact"];
	        this.Notes = source["Notes"];
	    }
	}
	export class ScopeRulesSettings {
	    AllowedHTTPMethods: number[];
	    DestructiveActionsForbid: boolean;
	    AuthBruteforceForbid: boolean;
	    DoSForbid: boolean;
	    CaptchaPolicy: number;
	    MaxRequestBodyKB: number;
	
	    static createFrom(source: any = {}) {
	        return new ScopeRulesSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AllowedHTTPMethods = source["AllowedHTTPMethods"];
	        this.DestructiveActionsForbid = source["DestructiveActionsForbid"];
	        this.AuthBruteforceForbid = source["AuthBruteforceForbid"];
	        this.DoSForbid = source["DoSForbid"];
	        this.CaptchaPolicy = source["CaptchaPolicy"];
	        this.MaxRequestBodyKB = source["MaxRequestBodyKB"];
	    }
	}
	export class Settings {
	    LLM: LLMSettings;
	    Program: ScopeProgramSettings;
	    Rules: ScopeRulesSettings;
	    Compliance: ComplianceSettings;
	    RateLimits: RateLimitSettings;
	    Assets: AssetSettings;
	    Authentication: AuthenticationSettings;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LLM = this.convertValues(source["LLM"], LLMSettings);
	        this.Program = this.convertValues(source["Program"], ScopeProgramSettings);
	        this.Rules = this.convertValues(source["Rules"], ScopeRulesSettings);
	        this.Compliance = this.convertValues(source["Compliance"], ComplianceSettings);
	        this.RateLimits = this.convertValues(source["RateLimits"], RateLimitSettings);
	        this.Assets = this.convertValues(source["Assets"], AssetSettings);
	        this.Authentication = this.convertValues(source["Authentication"], AuthenticationSettings);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

