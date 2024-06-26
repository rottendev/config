package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/joho/godotenv"
	"github.com/rottendev/config/pkg"
	"gopkg.in/yaml.v3"
)

type EnvProvider struct {
	Filename string
	prefix   string
}

func (e EnvProvider) Decode(data []byte, v interface{}) error {
	if e.Filename != "" {
		if err := godotenv.Load(e.Filename); err != nil {
			return err
		}
	}
	// Perform variable substitution on the YAML template
	configData := os.ExpandEnv(string(data))

	return yaml.Unmarshal([]byte(configData), v)
}

func (e EnvProvider) Encode(v any) ([]byte, error) {
	keys := make(map[string]interface{})
	pkg.GeneratePlaceholderMap(v, keys, e.prefix)

	// sort keys
	sortedKV := make([]string, 0, len(keys))
	for k := range keys {
		sortedKV = append(sortedKV, k)
	}
	sort.Strings(sortedKV)

	b := bytes.Buffer{}
	for _, k := range sortedKV {
		val := keys[k]
		// if val is pointer, get the value
		if reflect.ValueOf(val).Kind() == reflect.Ptr {
			val = reflect.ValueOf(val).Elem().Interface()
		}

		if reflect.ValueOf(val).Kind() == reflect.Slice {
			jsB, _ := json.Marshal(val)
			val = string(jsB)
		}
		_, _ = b.WriteString(fmt.Sprintf("%s=%v\n", k, val))
	}

	return b.Bytes(), nil
}
