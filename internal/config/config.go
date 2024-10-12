package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/n25a/BitCanary/internal/log"
	"go.uber.org/zap"
	"time"
)

// C is global config of BitCanary
var C *Config

type Config struct {
	HTTP            HTTP     `koanf:"http"`
	PrimaryAddress  string   `koanf:"primary_address"`
	CanaryAddress   string   `koanf:"canary_address"`
	SharedURLs      []string `koanf:"shared_urls"`
	Canary          Canary   `koanf:"canary"`
	UserIDHeaderKey string   `koanf:"user_id_header_key"`
	UserNestedKey   string   `koanf:"user_nested_key"`
}

type HTTP struct {
	Bind         string        `koanf:"bind"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
}

type Canary struct {
	Enabled   bool     `koanf:"enabled"`
	Bucket    Bucket   `koanf:"bucket"`
	Whitelist []uint64 `koanf:"whitelist"`
}

type Bucket struct {
	LowerBound uint32 `koanf:"lb"`
	UpperBound uint32 `koanf:"ub"`
}

// LoadConfig function will load the file located in path and return the parsed config.
// This function will panic on errors.
func LoadConfig(path string) *Config {
	// k is the global koanf instance. Use "." as the key path delimiter.
	k := koanf.New(".")

	// load default config in the beginning
	err := k.Load(structs.Provider(defaultConfig, "koanf"), nil)
	if err != nil {
		log.Logger.Fatal("error in loading the default config", zap.Error(err))
	}

	// load YAML config and merge into the previously loaded config.
	err = k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		log.Logger.Fatal("error in loading the config file", zap.Error(err))
	}

	var c Config
	err = k.Unmarshal("", &c)
	if err != nil {
		log.Logger.Fatal("error in unmarshalling the config file", zap.Error(err))
	}

	C = &c
	return &c
}
