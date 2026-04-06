package config

import (
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
)

const (
	ExifToolPath             = "exiftool.path"
	ImportCamVideoDefaultDst = "import.camvideo.default.targetDir"
	ImportGoProDefaultDst    = "import.gopro.default.targetDir"
	ImportSdPhotosDefaultDst = "import.sdPhotos.default.targetDir"
)

var defaults = map[string]map[string]string{
	"linux": map[string]string{
		//use OS preinstalled version
		ExifToolPath:             "",
		ImportCamVideoDefaultDst: "~/Media/camera",
		ImportGoProDefaultDst:    "~/Media/gopro",
		ImportSdPhotosDefaultDst: "~/Media/photos",
	},
	"windows": map[string]string{
		ExifToolPath:             "",
		ImportCamVideoDefaultDst: "d:\\video\\camera",
		ImportGoProDefaultDst:    "d:\\video\\gopro",
		ImportSdPhotosDefaultDst: "d:\\photos"},
}

type Config struct {
	ko *koanf.Koanf
}

func Empty() *Config {
	result := Config{ko: koanf.New(".")}
	return &result
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

	c.InitDefaults()
}

func (c *Config) SaveConfig(dstFile string) error {
	if err := os.MkdirAll(filepath.Dir(dstFile), 0o755); err != nil {
		return err
	}

	data, err := c.GetAsYaml()
	if err != nil {
		return err
	}

	return os.WriteFile(dstFile, data, 0o644)
}

func (c *Config) GetAsYaml() ([]byte, error) {
	return c.ko.Marshal(yaml.Parser())
}

func (c *Config) InitDefaults() {
	defs, ok := defaults[os.Getenv("OS")]
	if !ok {
		defs = defaults["linux"]
	}

	for key, value := range defs {
		c.SetDefault(key, value)
	}
}

func (c *Config) GetImportCamVideoDefaultDst() string {
	return c.ko.String(ImportCamVideoDefaultDst)
}

func (c *Config) GetExifToolPath() string {
	return c.ko.String(ExifToolPath)
}

func (c *Config) GetImportGoProDefaultDst() string {
	return c.ko.String(ImportGoProDefaultDst)
}

func (c *Config) GetImportSdPhotosDefaultDst() string {
	return c.ko.String(ImportSdPhotosDefaultDst)
}
