package config

import (
	"encoding/json"
	"os"

	"github.com/ferdiebergado/gopherkit/env"
)

type AppConfig struct {
	Env     string `json:"env,omitempty"`
	IsDebug bool   `json:"is_debug,omitempty"`
}

type DBConfig struct {
	User        string `json:"user,omitempty"`
	Pass        string `json:"pass,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	PingTimeout int    `json:"ping_timeout,omitempty"`
	DB          string `json:"db,omitempty"`
}

type ServerConfig struct {
	Port            int `json:"port,omitempty"`
	ShutdownTimeout int `json:"shutdown_timeout,omitempty"`
}

type Config struct {
	App    AppConfig    `json:"app,omitempty"`
	Db     DBConfig     `json:"db,omitempty"`
	Server ServerConfig `json:"server,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	config.App.Env = env.Get("ENV", config.App.Env)
	config.App.IsDebug = env.GetBool("DEBUG", config.App.IsDebug)

	config.Db.User = os.Getenv("POSTGRES_USER")
	config.Db.Pass = os.Getenv("POSTGRES_PASSWORD")
	config.Db.Host = env.Get("POSTGRES_HOST", config.Db.Host)
	config.Db.Port = env.GetInt("POSTGRES_PORT", config.Db.Port)
	config.Db.DB = env.Get("POSTGRES_DB", config.Db.DB)
	config.Db.PingTimeout = env.GetInt("DB_PING_TIMEOUT", config.Db.PingTimeout)

	config.Server.Port = env.GetInt("PORT", config.Server.Port)

	return &config, nil
}
