package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type GitAgentConfig struct {
	PollInterval      time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval    time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	HTTPClientTimeOut time.Duration `env:"CLIENT_TIMEOUT" envDefault:"20s"`
	ServerAddr        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

type GitServerConfig struct {
	ServerAddr string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

type Config struct {
	AgentConfig  GitAgentConfig
	ServerConfig GitServerConfig
	DebugMode    bool
}

func New() *Config {
	return &Config{
		AgentConfig: GitAgentConfig{
			PollInterval:      time.Duration(getEnvAsTime("POLL_INTERVAL", 2)),
			ReportInterval:    time.Duration(getEnvAsTime("REPORT_INTERVAL", 10)),
			HTTPClientTimeOut: time.Duration(getEnvAsTime("CLIENT_TIMEOUT", 20)),
			ServerAddr:        getEnv("ADDRESS", "127.0.0.1:8080"),
		},
		ServerConfig: GitServerConfig{
			ServerAddr: getEnv("ADDRESS", "127.0.0.1:8080"),
		},
		DebugMode: getEnvAsBool("DEBUG_MODE", true),
	}

}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsTime(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return int64(value)
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func (config *Config) ParseFlags() {
	flag.StringVar(&config.ServerConfig.ServerAddr, "a", config.ServerConfig.ServerAddr, "server address (host:port)")
	flag.DurationVar(&config.AgentConfig.ReportInterval, "r", config.AgentConfig.ReportInterval, "report interval (example: 10s)")
	flag.DurationVar(&config.AgentConfig.PollInterval, "p", config.AgentConfig.PollInterval, "poll interval (example: 10s)")
	flag.Parse()
}
