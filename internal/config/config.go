package config

import (
	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	Env   string `toml:"env"`
	Debug bool   `toml:"debug"`

	Database struct {
		Driver   string `toml:"driver"`
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Database string `toml:"database"`
		ShowSQL  bool   `toml:"show_sql"`
	} `toml:"database"`

	Web struct {
		ListenAddr string `toml:"listen_addr"`
	} `toml:"web"`
}

func LoadConfig(path string) (*AppConfig, error) {
	var config AppConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
