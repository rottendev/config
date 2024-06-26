package provider

import (
	"encoding/json"
)

type JSONProvider struct{}

func (JSONProvider) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (JSONProvider) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}
