package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Server Server `json:"server"`
}

type Server struct {
	RSAKey  RSAKey `json:"rsa_key"`
	Address string `json:"address"`
}

type RSAKey struct {
	E string `json:"e"`
	N string `json:"n"`
}

func MustLoad(configPath string) *Config {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to open config file: %s", configPath)
	}

	cfg := new(Config)
	err = json.Unmarshal(content, cfg)
	if err != nil {
		log.Fatalf("failed to read json: %s", err.Error())
	}

	return cfg
}
