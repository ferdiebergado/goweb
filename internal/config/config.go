package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	Env     string `json:"env,omitempty"`
	IsDebug bool   `json:"is_debug,omitempty"`
}

type DBConfig struct {
	User        string `json:"user,omitempty"`
	Pass        string `json:"pass,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        string `json:"port,omitempty"`
	PingTimeout int    `json:"ping_timeout,omitempty"`
}

type ServerConfig struct {
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

	return &config, nil
}
