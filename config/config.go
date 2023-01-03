package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type config struct {
	Mysql struct {
		Host string `yaml:"host"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		DB   string `yaml:"db"`
	}
	Redis struct {
		Host string `yaml:"host"`
		DB   int    `yaml:"db"`
		Pass string `yaml:"pass"`
	}
	Listen string `yaml:"listen"`
	Debug  bool   `yaml:"debug"`
}

var Global config

func Load() error {
	// Read config
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return fmt.Errorf("failed to read config: %s", err)
	}

	// Parse config to object
	err = yaml.Unmarshal(file, &Global)
	if err != nil {
		return fmt.Errorf("failed to parse config: %s", err)
	}

	return nil
}
