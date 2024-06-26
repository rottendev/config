package provider

import (
	"gopkg.in/yaml.v3"
)

type YamlProvider struct{}

func (YamlProvider) Decode(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (YamlProvider) Encode(v any) ([]byte, error) {
	return yaml.Marshal(v)
}
