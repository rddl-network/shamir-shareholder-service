package config

import (
	"sync"

	log "github.com/rddl-network/go-logger"
)

const DefaultConfigTemplate = `
service-host="{{ .ServiceHost }}"
service-port={{ .ServicePort }}
db-path="{{ .DBPath }}"
key-phrase="{{ .KeyPhrase }}"
certs-path="{{ .CertsPath }}"
log-level="{{ .LogLevel }}"
`

type Config struct {
	ServicePort int    `mapstructure:"service-port"`
	ServiceHost string `mapstructure:"service-host"`
	DBPath      string `mapstructure:"db-path"`
	KeyPhrase   string `mapstructure:"key-phrase"`
	CertsPath   string `mapstructure:"certs-path"`
	LogLevel    string `mapstructure:"log-level"`
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
		CertsPath:   "./certs/",
		LogLevel:    log.ERROR,
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
