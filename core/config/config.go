package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds the Agent configuration
type Config struct {
	// todo: add development URL
	AccountID  string `yaml:"account_id"`
	LicenseKey string `yaml:"license_key"`
}

// Load parses a Config object from a YAML file whose
func Load(yamlFilePath string) (Config, error) {
	var cfg Config
	yml, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(yml, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
