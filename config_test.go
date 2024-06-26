package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testConfig struct {
	App struct {
		Name string `yaml:"name" json:"name" default:"app"`
		Port int    `yaml:"port" default:"8080"`
	} `json:"app" toml:"sSs"`
	Region   string   `yaml:"region"  json:"region" default:"us-west-1"`
	FilesDir *string  `yaml:"files_dir" json:"filesDir"`
	Modules  []string `yaml:"modules" json:"modulesAA" default:"[\"module1\", \"module2\"]"`
}

func TestConfig_LoadConfig(t *testing.T) {
	type args struct {
		conf     interface{}
		filename string
		data     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Yaml Success",
			args: args{
				conf:     &testConfig{},
				filename: "testdata/config.test.yaml",
				data:     []byte("app:\n    name: appYaml\n    port: 8081\nregion: us-west-2\nfiles_dir: null\nmodules:\n    - module1\n    - module2\n"),
			},
		},
		{
			name: "JSON Success",
			args: args{
				conf:     &testConfig{},
				filename: "testdata/config.test.json",
				data:     []byte("{\"app\":{\"name\":\"appJson\",\"Port\":8082},\"region\":\"us-west-3\",\"filesDir\":null,\"modulesAA\":[\"module1\",\"module2\"]}"),
			},
		},
		{
			name: "Toml Success",
			args: args{
				conf:     &testConfig{},
				filename: "testdata/config.test.toml",
				data:     []byte("Region = \"us-west-3\"\nFilesDir = \"appToml\"\nModules = [\"module4\", \"module5\"]\n\n[sSs]\n  Name = \"appToml\"\n  Port = 8084\n"),
			},
		},
		{
			name: "Env Success",
			args: args{
				conf:     &testConfig{},
				filename: "testdata/config.test.env",
				data:     []byte("app:\n    name: ${APP_NAME}\n    port: ${APP_PORT}\nfiles_dir: ${FILES_DIR}\nmodules: ${MODULES}\nregion: ${REGION}\n"),
			},
		},
		{
			name: "File Not Found",
			args: args{
				conf:     &testConfig{},
				filename: "testdata/config.test.a.txt",
				data:     []byte("app:\n    name: appYaml\n    port: 8081\nregion: us-west-2\nfiles_dir: null\nmodules:\n    - module1\n    - module2\n"),
			},
			wantErr: true,
		},
		{
			name: "Invalid Config",
			args: args{
				conf:     &testConfig{},
				filename: "testdata/config.test.env.yaml",
				data:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.conf, tt.args.filename, tt.args.data); (err != nil) != tt.wantErr {
				require.Failf(t, "", "LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			cfg := GetConfig()
			require.NotNil(t, cfg)
			require.NotNil(t, cfg.parsedConfig)
			require.Equal(t, DetectConfigType(tt.args.filename), cfg.cfgType)
			require.Equal(t, tt.args.conf, cfg.parsedConfig)
			data, err := cfg.Encode()
			require.NoError(t, err)

			if tt.name == "Env Success" {
				require.Len(t, tt.args.conf.(*testConfig).Modules, 2)
				require.Equal(t, "APP_NAME=appEnv\nAPP_PORT=8085\nFILES_DIR=appEnv\nMODULES=[\"module6\",\"module7\"]\nREGION=us-west-4\n", string(data))
				return
			}
			require.Equal(t, tt.args.data, data, fmt.Sprintf("expected: %s, got: %s", string(tt.args.data), string(data)))
		})
	}
}

func resetEnv() {
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_PORT")
	_ = os.Unsetenv("FILES_DIR")
	_ = os.Unsetenv("MODULES")
	_ = os.Unsetenv("REGION")
}

func TestConfig_ENV(t *testing.T) {
	resetEnv()

	setting := testConfig{}
	err := LoadConfig(&setting, "", nil)
	require.Error(t, err)
	require.Equal(t, "missing template data", err.Error())

	err = LoadConfig(&setting, "", []byte("non-env data"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode yaml:")

	data, err := os.ReadFile("testdata/config.test.env.yaml")
	require.NoError(t, err)

	t.Run("default env", func(t *testing.T) {
		defer resetEnv()
		cfg, _ := New(EnvConfig, "")
		err = cfg.LoadConfig(&setting, data)
		require.NoError(t, err)
		require.Equal(t, "app", setting.App.Name)
		require.Equal(t, 8080, setting.App.Port)
		require.Equal(t, "us-west-1", setting.Region)
		require.Nil(t, setting.FilesDir)
		require.Len(t, setting.Modules, 0)
	})

	t.Run("Machine env", func(t *testing.T) {
		defer resetEnv()
		_ = os.Setenv("APP_NAME", "appEnv1")
		_ = os.Setenv("APP_PORT", "9085")
		_ = os.Setenv("FILES_DIR", "1111")
		_ = os.Setenv("MODULES", "[\"22\",\"33\"]")
		_ = os.Setenv("REGION", "us-west-4")

		cfg, _ := New(EnvConfig, "")
		err = cfg.LoadConfig(&setting, data)
		require.NoError(t, err)
		require.Equal(t, "appEnv1", setting.App.Name)
		require.Equal(t, 9085, setting.App.Port)
		require.Equal(t, "us-west-4", setting.Region)
		require.Equal(t, "1111", *setting.FilesDir)
		require.Len(t, setting.Modules, 2)
		require.Equal(t, "22", setting.Modules[0])
		require.Equal(t, "33", setting.Modules[1])
	})

	t.Run("File env", func(t *testing.T) {
		defer resetEnv()

		cfg, _ := WithFile("testdata/config.test.env")
		err = cfg.LoadConfig(&setting, data)

		err = LoadConfig(&setting, "testdata/config.test.env", data)
		require.NoError(t, err)
		require.Equal(t, "appEnv", setting.App.Name)
		require.Equal(t, 8085, setting.App.Port)
		require.Equal(t, "us-west-4", setting.Region)
		require.Equal(t, "appEnv", *setting.FilesDir)
		require.Len(t, setting.Modules, 2)
		require.Equal(t, "module6", setting.Modules[0])
		require.Equal(t, "module7", setting.Modules[1])
	})

	t.Run("Machine and file env conflicts", func(t *testing.T) {
		defer resetEnv()
		_ = os.Setenv("APP_NAME", "App Env from machine")
		_ = os.Setenv("FILES_DIR", "")
		_ = os.Setenv("REGION", "GOOGLE")

		cfg, _ := WithFile("testdata/config.test.env")
		err = cfg.LoadConfig(&setting, data)
		require.NoError(t, err)

		require.Equal(t, "App Env from machine", setting.App.Name)
		require.Equal(t, "GOOGLE", setting.Region)
		require.Nil(t, setting.FilesDir)

		require.Equal(t, 8085, setting.App.Port)
		require.Len(t, setting.Modules, 2)
		require.Equal(t, "module6", setting.Modules[0])
		require.Equal(t, "module7", setting.Modules[1])
	})

	t.Run("Invalid env", func(t *testing.T) {
		defer resetEnv()
		_ = os.Setenv("APP_PORT", "invalid")

		cfg, _ := WithFile("testdata/config.test.env")
		err = cfg.LoadConfig(&setting, data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "decode yaml:")
	})

	t.Run("static file", func(t *testing.T) {
		defer resetEnv()
		_ = os.Setenv("APP_NAME", "App Env from machine")

		cfg, _ := WithFile("testdata/config.test.env")
		err = cfg.LoadConfig(&setting, data)

		err = LoadConfig(&setting, "testdata/config.test.env", []byte("app:\n    name: StaticName\n    port: ${APP_PORT}\nfiles_dir: ${FILES_DIR}\nmodules: ${MODULES}\nregion: ${REGION}\n"))
		require.NoError(t, err)
		require.Equal(t, "StaticName", setting.App.Name)
	})
}

func TestErrUnsupportedConfigType_Error(t *testing.T) {
	unErr := ErrUnsupportedConfigType("test")
	require.Equal(t, "unsupported config type: \"test\"", unErr.Error())

	// It will default to env
	_, err := WithFile("testdata/config.test.txt")
	require.NoError(t, err)

	// This one should fail
	_, err = New("txt", "testdata/config.test.txt")
	require.Error(t, err)
}
