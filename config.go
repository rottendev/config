package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/rottendev/config/provider"
)

// Provider is the interface that wraps the basic methods for a configuration provider.
type Provider interface {
	Decode(data []byte, v interface{}) error
	Encode(v any) ([]byte, error)
}

type Config struct {
	cfgType      Type              // The configuration type.
	providers    map[Type]Provider // The configuration providers.
	filename     string            // The configuration filename.
	parsedConfig interface{}
}

var c *Config

func WithFile(filename string) (*Config, error) {
	cfgType := DetectConfigType(filename)
	return New(cfgType, filename)
}

func New(cfgType Type, filename string) (*Config, error) {
	if filename != "" {
		if _, err := os.Stat(filename); err != nil {
			return nil, fmt.Errorf("config %w", err)
		}
	}
	c = &Config{
		cfgType:  cfgType,
		filename: filename,
	}

	return c.initProviders()
}

// initProviders initializes the configuration providers.
func (c *Config) initProviders() (*Config, error) {
	c.providers = make(map[Type]Provider)

	switch c.cfgType {
	case JSONConfig:
		c.providers[JSONConfig] = &provider.JSONProvider{}
	case YamlConfig:
		c.providers[YamlConfig] = &provider.YamlProvider{}
	case TomlConfig:
		c.providers[TomlConfig] = &provider.TomlProvider{}
	case EnvConfig:
		c.providers[EnvConfig] = &provider.EnvProvider{Filename: c.filename}
	default:
		return nil, ErrUnsupportedConfigType(c.cfgType)
	}

	return c, nil
}

func (c *Config) getProvider() (Provider, error) {
	p, ok := c.providers[c.cfgType]
	if !ok {
		return nil, ErrUnsupportedConfigType(c.cfgType)
	}

	return p, nil
}

// LoadConfig loads the configuration from the providers.
func LoadConfig(conf interface{}, filename string, data []byte) error {
	var err error
	if c == nil || c.filename != filename {
		c, err = WithFile(filename)
		if err != nil {
			return err
		}
	}

	return c.LoadConfig(conf, data)
}
func (c *Config) LoadConfig(conf interface{}, data []byte) error {
	if c.cfgType == EnvConfig {
		if data == nil {
			return fmt.Errorf("missing template data")
		}
	} else {
		var err error
		data, err = os.ReadFile(c.filename)
		if err != nil {
			return err
		}
	}

	p, err := c.getProvider()
	if err != nil {
		return err
	}

	if err = defaults.Set(conf); err != nil {
		return err
	}

	if err = p.Decode(data, conf); err != nil {
		return fmt.Errorf("decode %w", err)
	}

	c.parsedConfig = conf

	return nil
}

// GetConfig returns the configuration.
func GetConfig() *Config {
	return c
}

func (c *Config) Encode() ([]byte, error) {
	p, err := c.getProvider()
	if err != nil {
		return nil, err
	}

	return p.Encode(c.parsedConfig)
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{type: %s, filename: %s}", c.cfgType, c.filename)
}
