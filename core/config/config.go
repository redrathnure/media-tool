package config

import (
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	ko             *koanf.Koanf
}

func Empty() *Config {
	return &Config{ko: koanf.New(".")}
}

func (c *Config) GetString(key string) string {
	return c.ko.String(key)
}

func (c *Config) SetDefault(key string, value interface{}) {
	if !c.ko.Exists(key) {
		_ = c.ko.Set(key, value)
	}
}

func (c *Config) Set(key string, value interface{}) error {
	return c.ko.Set(key, value)
}

func (c *Config) Reset() {
	c.ko = koanf.New(".")
}

func (c *Config) WriteConfigAs(dstFile string) error {
	if err := os.MkdirAll(filepath.Dir(dstFile), 0o755); err != nil {
		return err
	}

	data, err := c.ko.Marshal(yaml.Parser())
	if err != nil {
		return err
	}

	return os.WriteFile(dstFile, data, 0o644)
}
