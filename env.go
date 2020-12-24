package slice

import (
	"os"
	"strings"
)

func parseEnv(env string) Env {
	return Env(strings.ToLower(env))
}

var lookupEnv = os.LookupEnv

// Env
type Env string

// String converts environment value to predefined strings.
func (e Env) String() string {
	return string(e)
}
