package coffy

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func ParseFile(file *os.File) (*Config, error) {
	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	return Parse(content)
}

func Parse(cfg string) (*Config, error) {
	var data = &Config{}
	err := yaml.Unmarshal([]byte(cfg), data)
	if err != nil {
		return nil, err
	}
	err = validateConfig(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func validateConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	return validateServer(cfg.Server)
}

func validateServer(s *ServerCfg) error {
	if s == nil {
		return MissingPropertyError{"server", "missing property"}
	}
	if s.Port == 0 {
		return MissingPropertyError{"port", "missing property"}
	}
	return nil
}

type Config struct {
	Server *ServerCfg `yaml:"server"`
}

type ServerCfg struct {
	Port int `yaml:"port"`
}

type MissingPropertyError struct {
	Property string
	Message  string
}

func (e MissingPropertyError) Error() string {
	return fmt.Sprintf("Property: '%s'. %s", e.Property, e.Message)
}
