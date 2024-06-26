package provider

import (
	"github.com/BurntSushi/toml"
)

type TomlProvider struct{}

func (TomlProvider) Decode(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func (TomlProvider) Encode(v any) ([]byte, error) {
	return toml.Marshal(v)
}
