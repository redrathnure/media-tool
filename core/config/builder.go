package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const configFileName = "media-tool.yml"

type ConfigBuilder struct {
	ko          *koanf.Koanf
	configFiles []string
}

func NewConfig() *ConfigBuilder {
	return &ConfigBuilder{ko: koanf.New(".")}
}

func (c *ConfigBuilder) AddConfigLocation(pathComponents ...string) {
	path := path.Join(pathComponents...)
	path = c.expandPath(path)
	candidate := filepath.Join(path, configFileName)
	c.configFiles = append(c.configFiles, candidate)
}

func (c *ConfigBuilder) SetConfigFile(confFile string) {
	c.configFiles = append(c.configFiles, c.expandPath(confFile))
}

func (c *ConfigBuilder) AutomaticEnv() {
	_ = c.ko.Load(env.ProviderWithValue("", ".", func(key string, value string) (string, interface{}) {
		transformed := strings.ToLower(strings.TrimPrefix(key, "MEDIA_TOOL_"))
		transformed = strings.ReplaceAll(transformed, "_", ".")
		if transformed == "" {
			return "", nil
		}
		return transformed, value
	}), nil)
}

func (c *ConfigBuilder) Build() (*Config, error) {
	c.AutomaticEnv()

	for _, confFile := range c.configFiles {
		if err := c.loadFile(confFile); err != nil {
			return nil, err
		}

	}

	result := Config{ko: c.ko}
	result.InitDefaults()

	return &result, nil
}

func (c *ConfigBuilder) loadFile(confFile string) error {
	if confFile == "" {
		return fmt.Errorf("config file was not specified")
	}
	// Check file exists and it is file
	info, err := os.Stat(confFile)
	if err != nil || info.IsDir() {
		return nil
	}

	log.Debugf("Loading configuration from '%s' file...", confFile)
	if err := c.ko.Load(file.Provider(confFile), yaml.Parser()); err != nil {
		return fmt.Errorf("Unable to load '%v' config: %v", confFile, err)
	}

	return nil
}

func (c *ConfigBuilder) expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}

	return filepath.Clean(os.ExpandEnv(path))
}
