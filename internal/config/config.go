package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ferdiebergado/gopherkit/env"
)

type EnvConfig struct {
	Env     string `json:"env,omitempty"`
	IsDebug bool   `json:"is_debug,omitempty"`
}

type DBConfig struct {
	Driver          string `json:"driver,omitempty"`
	User            string `json:"user,omitempty"`
	Pass            string `json:"pass,omitempty"`
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	SSLMode         string `json:"ssl_mode,omitempty"`
	PingTimeout     int    `json:"ping_timeout,omitempty"`
	DB              string `json:"db,omitempty"`
	MaxOpenConns    int    `json:"max_open_conns,omitempty"`
	MaxIdleConns    int    `json:"max_idle_conns,omitempty"`
	ConnMaxIdle     int    `json:"conn_max_idle,omitempty"`
	ConnMaxLifetime int    `json:"conn_max_lifetime,omitempty"`
}

type ServerConfig struct {
	Port            int `json:"port,omitempty"`
	ReadTimeout     int `json:"read_timeout,omitempty"`
	WriteTimeout    int `json:"write_timeout,omitempty"`
	IdleTimeout     int `json:"idle_timeout,omitempty"`
	ShutdownTimeout int `json:"shutdown_timeout,omitempty"`
}

type TemplateConfig struct {
	Path         string `json:"path,omitempty"`
	LayoutFile   string `json:"layout_file,omitempty"`
	PartialsPath string `json:"partials_path,omitempty"`
	PagesPath    string `json:"pages_path,omitempty"`
}

type Config struct {
	App      EnvConfig      `json:"app,omitempty"`
	Db       DBConfig       `json:"db,omitempty"`
	Server   ServerConfig   `json:"server,omitempty"`
	Template TemplateConfig `json:"template,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	path = filepath.Clean(path)
	configFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file %s: %w", path, err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	config.App.Env = env.Get("ENV", config.App.Env)
	config.App.IsDebug = env.GetBool("DEBUG", config.App.IsDebug)

	config.Db.User = env.MustGet("POSTGRES_USER")
	config.Db.Pass = env.MustGet("POSTGRES_PASSWORD")
	config.Db.Host = env.Get("POSTGRES_HOST", config.Db.Host)
	config.Db.Port = env.GetInt("POSTGRES_PORT", config.Db.Port)
	config.Db.DB = env.Get("POSTGRES_DB", config.Db.DB)
	config.Db.SSLMode = env.MustGet("POSTGRES_SSLMODE")
	config.Db.PingTimeout = env.GetInt("DB_PING_TIMEOUT", config.Db.PingTimeout)
	config.Db.MaxOpenConns = env.GetInt("DB_MAX_OPEN_CONNS", config.Db.MaxOpenConns)
	config.Db.MaxIdleConns = env.GetInt("DB_MAX_IDLE_CONNS", config.Db.MaxIdleConns)
	config.Db.ConnMaxLifetime = env.GetInt("DB_CONN_MAX_LIFETIME", config.Db.ConnMaxLifetime)
	config.Db.ConnMaxIdle = env.GetInt("DB_CONN_MAX_IDLE", config.Db.ConnMaxIdle)

	config.Server.Port = env.GetInt("PORT", config.Server.Port)

	return &config, nil
}
