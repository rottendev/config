package pkg

import (
	"fmt"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "simple",
			in:   "Simple",
			want: "simple",
		},
		{
			name: "camelCase",
			in:   "camelCase",
			want: "camel_case",
		},
		{
			name: "HTTPServer",
			in:   "HTTPServer",
			want: "http_server",
		},
		{
			name: "HTTPServerURL",
			in:   "HTTPServerURL",
			want: "http_server_url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.in); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePlaceholderMap(t *testing.T) {
	type testInnerStruct struct {
		URL string      `yaml:"url"`
		B   interface{} `yaml:"b"`
	}
	type testStruct struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Inner []testInnerStruct
	}
	keys := make(map[string]interface{})
	prefix := "APP_"
	test := testStruct{
		Host: "localhost",
		Port: 8080,
		Inner: []testInnerStruct{
			{URL: "http://localhost:8080"},
			{URL: "http://localhost:8081", B: 1},
		},
	}

	got := GeneratePlaceholderMap(&test, keys, prefix)
	fmt.Println(got)
	want := map[string]interface{}{
		"host": "${APP_HOST}",
		"port": "${APP_PORT}",
		"inner": []map[string]interface{}{
			{"url": "${APP_INNER_0_URL}", "b": "${APP_INNER_0_B}"},
			{"url": "${APP_INNER_1_URL}", "b": "${APP_INNER_1_B}"},
		},
	}
	if len(got) != len(want) {
		t.Errorf("GeneratePlaceholderMap() = %v, want %v", got, want)
	}
}
