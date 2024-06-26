# Config

## About
Currently, it supports following configuration formats:

yaml, yml, json, toml, env.

```go
package main

import (
	"fmt"
	"log"
)

type PersonConfig struct {
    Name    string
    Age     int
    Weight  float64
    Address struct {
        City    string
        Country string
    } `json:"address,omitempty"`
}

func main() {
    var cfg PersonConfig
    err := config.LoadConfig(&cfg, "./etc/project/config.yaml", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(cfg)
}
```