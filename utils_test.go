package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectConfigType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     Type
	}{
		{
			name:     "json",
			filename: "config.json",
			want:     JSONConfig,
		},
		{
			name:     "yaml",
			filename: "config.yaml",
			want:     YamlConfig,
		},
		{
			name:     "yml",
			filename: "config.yml",
			want:     YamlConfig,
		},
		{
			name:     "toml",
			filename: "config.toml",
			want:     TomlConfig,
		},
		{
			name:     "env",
			filename: "config.env",
			want:     EnvConfig,
		},
		{
			name:     "unknown",
			filename: "config.unknown",
			want:     EnvConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectConfigType(tt.filename); got != tt.want {
				t.Errorf("DetectConfigType() = %v, want %v", got, tt.want)
			}
		})
	}
}

const ymlContent = `app:
    name: app
    port: 8080
region: us-west-1
files_dir: null
modules:
    - module1
    - module2
`

const jsonContent = `{"App":{"Name":"app","Port":8080},"Region":"us-west-1","FilesDir":null,"Modules":["module1","module2"]}
`

const tomlContent = `Region = "us-west-1"
Modules = ["module1", "module2"]

[App]
  Name = "app"
  Port = 8080
`
const envYamlTemplate = `app:
    name: ${APP_NAME}
    port: ${APP_PORT}
files_dir: ${FILES_DIR}
modules: ${MODULES}
region: ${REGION}
`
const envContent = `APP_NAME=app
APP_PORT=8080
FILES_DIR=<nil>
MODULES=[module1 module2]
REGION=us-west-1
`

func TestExportStructs(t *testing.T) {
	type Conf struct {
		App struct {
			Name string `yaml:"name" default:"app"`
			Port int    `yaml:"port" default:"8080"`
		}
		Region   string   `yaml:"region" default:"us-west-1"`
		FilesDir *string  `yaml:"files_dir"`
		Modules  []string `yaml:"modules" default:"[\"module1\", \"module2\"]"`
	}

	tests := []struct {
		name    string
		output  string
		cfgType Type
		want    string
	}{
		{
			name:    "yaml",
			output:  "config.test.yaml",
			cfgType: YamlConfig,
			want:    ymlContent,
		},
		{
			name:    "json",
			output:  "config.test.json",
			cfgType: JSONConfig,
			want:    jsonContent,
		},
		{
			name:    "toml",
			output:  "config.test.toml",
			cfgType: TomlConfig,
			want:    tomlContent,
		},
		{
			name:    "env",
			output:  "config.test.env.yaml",
			cfgType: EnvConfig,
			want:    envYamlTemplate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Conf{}
			fileName := ExportStructs(&cfg, tt.cfgType, tt.output)
			defer func() {
				_ = os.Remove(fileName)
				if tt.cfgType == EnvConfig {
					_ = os.Remove("config.sample.env")
				}
			}()
			require.Equal(t, tt.output, fileName)

			file, err := os.ReadFile(fileName)
			if err != nil {
				t.Error(err)
			}
			require.NotEmpty(t, file)
			require.Equal(t, tt.want, string(file))

			if tt.cfgType == EnvConfig {
				file, err = os.ReadFile("config.sample.env")
				if err != nil {
					t.Error(err)
				}
				require.NotEmpty(t, file)
				require.Equal(t, envContent, string(file))
			}
		})
	}
}
