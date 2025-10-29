package settings

type LLMSettings struct {
	Provider           string
	AllowInternetAccess bool
	MaxContextTokens    int
}

type Settings struct {
	LLM LLMSettings
}

func Default() Settings {
	return Settings{
		LLM: LLMSettings{
			Provider:           "local",
			AllowInternetAccess: false,
			MaxContextTokens:    32000,
		},
	}
}