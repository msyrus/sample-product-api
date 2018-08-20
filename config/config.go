package config

import (
	"io"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Application holds application configuration
type Application struct {
	GracefulWait time.Duration `yaml:"gracefulWait"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Postgres     Postgres      `yaml:"postgres"`
}

// Postgres holds postgres configuration
type Postgres struct {
	URI string `yml:"uri"`
}

// Parse return Application configuration from reader r
// GracefulWait, ReadTimeout, WriteTimeout, IdleTimeout are read as second
func Parse(r io.Reader) (Application, error) {
	cfg := Application{}
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return Application{}, err
	}
	app := Application{
		GracefulWait: cfg.GracefulWait * time.Second,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
		IdleTimeout:  cfg.IdleTimeout * time.Second,
		Host:         cfg.Host,
		Port:         cfg.Port,
		Postgres:     cfg.Postgres,
	}
	return app, nil
}
