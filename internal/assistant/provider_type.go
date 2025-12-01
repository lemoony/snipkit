package assistant

// ProviderType represents the type of LLM provider.
type ProviderType string

const (
	ProviderTypeOpenAI           = ProviderType("openai")
	ProviderTypeAnthropic        = ProviderType("anthropic")
	ProviderTypeGemini           = ProviderType("gemini")
	ProviderTypeOllama           = ProviderType("ollama")
	ProviderTypeOpenAICompatible = ProviderType("openai-compatible")
)

// SupportedProviders lists all available provider types.
var SupportedProviders = []ProviderType{
	ProviderTypeOpenAI,
	ProviderTypeAnthropic,
	ProviderTypeGemini,
	ProviderTypeOllama,
	ProviderTypeOpenAICompatible,
}

// ProviderInfo contains display information for a provider.
type ProviderInfo struct {
	Name             string
	Description      string
	DefaultModel     string
	DefaultAPIKeyEnv string
}

// GetProviderInfo returns display information for a provider type.
func GetProviderInfo(providerType ProviderType) ProviderInfo {
	info := map[ProviderType]ProviderInfo{
		ProviderTypeOpenAI: {
			Name:             "OpenAI",
			Description:      "Use OpenAI (GPT-4, GPT-4o) as assistant",
			DefaultModel:     "gpt-4o",
			DefaultAPIKeyEnv: "SNIPKIT_OPENAI_API_KEY",
		},
		ProviderTypeAnthropic: {
			Name:             "Anthropic",
			Description:      "Use Anthropic Claude as assistant",
			DefaultModel:     "claude-sonnet-4-20250514",
			DefaultAPIKeyEnv: "SNIPKIT_ANTHROPIC_API_KEY",
		},
		ProviderTypeGemini: {
			Name:             "Google Gemini",
			Description:      "Use Google Gemini as assistant",
			DefaultModel:     "gemini-1.5-flash",
			DefaultAPIKeyEnv: "SNIPKIT_GEMINI_API_KEY",
		},
		ProviderTypeOllama: {
			Name:             "Ollama",
			Description:      "Use local Ollama models as assistant",
			DefaultModel:     "llama3",
			DefaultAPIKeyEnv: "",
		},
		ProviderTypeOpenAICompatible: {
			Name:             "OpenAI-Compatible",
			Description:      "Use any OpenAI-compatible API (Together.ai, Groq, etc.)",
			DefaultModel:     "",
			DefaultAPIKeyEnv: "",
		},
	}

	if i, ok := info[providerType]; ok {
		return i
	}
	return ProviderInfo{}
}
