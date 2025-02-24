package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rottendev/config/pkg"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type Type string

const (
	YamlConfig Type = "yaml"
	JSONConfig Type = "json"
	TomlConfig Type = "toml"
	EnvConfig  Type = "env"
)

// DetectConfigType detects the type of configuration file based on its extension.
func DetectConfigType(filename string) Type {
	extension := strings.ToLower(filepath.Ext(filename))

	switch extension {
	case ".json":
		return JSONConfig
	case ".yaml", ".yml":
		return YamlConfig
	case ".toml":
		return TomlConfig
	case ".env":
		return EnvConfig
	default:
		return EnvConfig
	}
}

func outFileName(output string, cfgType Type) string {
	if output == "" {
		switch cfgType {
		case JSONConfig:
			return "config.sample.json"
		case YamlConfig:
			return "config.sample.yaml"
		case TomlConfig:
			return "config.sample.toml"
		default:
			return "config.env.yaml"
		}
	}
	return output
}

func ExportStructs(structure interface{}, cfgType Type, output string) string {
	dErr := defaults.Set(structure)
	if dErr != nil {
		panic(dErr)
	}

	outputFile := outFileName(output, cfgType)
	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	if cfgType == YamlConfig {
		if err = yaml.NewEncoder(f).Encode(structure); err != nil {
			panic(err)
		}
	}
	if cfgType == JSONConfig {
		if err = json.NewEncoder(f).Encode(structure); err != nil {
			panic(err)
		}
	}
	if cfgType == TomlConfig {
		if err = toml.NewEncoder(f).Encode(structure); err != nil {
			panic(err)
		}
	}
	if cfgType == EnvConfig {
		keys := make(map[string]interface{})
		placeholderMap := pkg.GeneratePlaceholderMap(structure, keys, "")

		var f1 *os.File
		// Write envs to .env.dev file
		f1, err = os.Create("config.sample.env")
		if err != nil {
			panic(err)
		}

		// sort keys
		sortedKV := make([]string, 0, len(keys))
		for k := range keys {
			sortedKV = append(sortedKV, k)
		}
		sort.Strings(sortedKV)

		for _, k := range sortedKV {
			v := keys[k]
			_, _ = f1.WriteString(fmt.Sprintf("%s=%v\n", k, v))
		}
		_ = f1.Close()

		if err = yaml.NewEncoder(f).Encode(placeholderMap); err != nil {
			panic(err)
		}
	}

	return outputFile
}
