package config

import (
	"fmt"
)

type ErrUnsupportedConfigType string

func (e ErrUnsupportedConfigType) Error() string {
	return fmt.Sprintf("unsupported config type: %q", string(e))
}
