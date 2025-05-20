package config

// holds application configuration
type AppConfig struct {
	ServerPort int
}

func NewAppConfig() *AppConfig {
	cfg := &AppConfig{
		ServerPort: 8080,
	}

	return cfg
}
