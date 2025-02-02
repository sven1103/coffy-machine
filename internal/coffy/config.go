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

	if err := validateServer(cfg.Server); err != nil {
		return err
	}

	if err := validateDatabase(cfg.Database); err != nil {
		return err
	}
	return nil
}

func validateDatabase(c *DbCfg) error {
	if c == nil {
		return MissingPropertyError{"database", "missing property"}
	}
	if c.Path == "" {
		return MissingPropertyError{"path", "missing property"}
	}
	return nil
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
	Server   *ServerCfg `yaml:"server"`
	Database *DbCfg     `yaml:"database"`
}

type ServerCfg struct {
	Port int `yaml:"port"`
}

type DbCfg struct {
	Path string `yaml:"path"`
}

type MissingPropertyError struct {
	Property string
	Message  string
}

func (e MissingPropertyError) Error() string {
	return fmt.Sprintf("%s: '%s'", e.Message, e.Property)
}
