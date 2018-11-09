package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"time"
)

const (
	defaultHarvestRate = 5 * time.Second
)

type Config struct {
	// todo: add development URL
	AccountId string `yaml:"account_id"`
	LicenseKey string `yaml:"license_key"`
}

func Load(file string) (Config, error) {
	var cfg Config
	yml, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(yml, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
