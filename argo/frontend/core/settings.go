package settings

struct LLMSettings {
	provider: string = "local"
	allow_internet_access: bool = "false"
	max_context_tokens: int = 32000
}

struct Settings {
	llm_settings: LLMSettings
}