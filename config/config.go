package config

import "sync"

const DefaultConfigTemplate = `
service-host="{{ .ServiceHost }}"
service-port={{ .ServicePort }}
db-path="{{ .DBPath }}"
key-phrase="{{ .KeyPhrase }}"
`

type Config struct {
	ServicePort int    `mapstructure:"service-port"`
	ServiceHost string `mapstructure:"service-host"`
	DBPath      string `mapstructure:"db-path"`
	KeyPhrase   string `mapstructure:"key-phrase"`
}

// global singletonjj
var (
	config     *Config
	initConfig sync.Once
)

// DefaultConfig returns RDDL-2-PLMNT default config
func DefaultConfig() *Config {
	return &Config{
		ServicePort: 8080,
		ServiceHost: "localhost",
		DBPath:      "./data",
		KeyPhrase:   "keyphrase",
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
