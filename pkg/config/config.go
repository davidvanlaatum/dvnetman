package config

import (
	"dvnetman/pkg/utils"
	"encoding/base64"
	"net/url"
	"os"
)
import "github.com/pkg/errors"
import "gopkg.in/yaml.v3"

type ListenAddress struct {
	Addr string
}

type MongoConfig struct {
	URL      string
	Database string
}

type OpenIDConnectAuthConfig struct {
	ClientID         string   `yaml:"clientID"`
	ClientSecret     string   `yaml:"clientSecret"`
	Provider         string   `yaml:"provider"`
	AutoDiscoveryURL string   `yaml:"autoDiscoveryURL"`
	Scopes           []string `yaml:"scopes"`
}

type AuthConfig struct {
	OpenIDConnect *OpenIDConnectAuthConfig `yaml:"openIDConnect"`
}

type SessionConfig struct {
	HashKey  string
	BlockKey string
}

func base64ToBytes(base64Str string) []byte {
	if base64Str == "" {
		return nil
	}
	if t, err := base64.StdEncoding.DecodeString(base64Str); err != nil {
		panic(err)
	} else {
		return t
	}
}

func (s *SessionConfig) HashKeyBytes() []byte {
	return base64ToBytes(s.HashKey)
}

func (s *SessionConfig) BlockKeyBytes() []byte {
	return base64ToBytes(s.BlockKey)
}

type Config struct {
	URL     string
	Listen  []ListenAddress
	Mongo   MongoConfig
	Auth    []AuthConfig
	Session SessionConfig
	url     *url.URL
}

func (c *Config) GetURL() (_ *url.URL, err error) {
	if c.url != nil {
		return c.url, nil
	}
	c.url, err = url.Parse(c.URL)
	return c.url, err
}

func LoadConfig(path string) (config *Config, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, errors.WithMessagef(err, "failed to open config file %s", path)
	}
	defer utils.PropagateError(f.Close, &err, "failed to close config file")
	err = errors.WithMessagef(yaml.NewDecoder(f).Decode(&config), "failed to decode config file %s", path)
	if config != nil {
		if len(config.Listen) == 0 {
			config.Listen = []ListenAddress{{Addr: ":8080"}}
		}
	}
	return
}
