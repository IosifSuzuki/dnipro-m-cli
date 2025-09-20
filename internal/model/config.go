package model

import (
	"dniprom-cli/pkg/logger"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ProductCodes      []string `yaml:"product_codes"`
	BaseURL           string   `yaml:"base_url"`
	ENV               string   `yaml:"env"`
	FileID            string   `yaml:"file_id"`
	GoogleCredentials string   `yaml:"google_credentials"`
}

func LoadConfig() (*Config, error) {
	path := "./config.yml"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func (c *Config) GetLoggerENV() logger.ENV {
	env, _ := logger.ENVFromString(c.ENV)
	return env
}
