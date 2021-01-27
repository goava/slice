package slice

import (
	"os"
	"strings"
)

const (
	envDev  string = "dev"
	envProd string = "prod"
)

func parseEnv(env string) Env {
	if env == "" {
		return Env(envProd)
	}
	return Env(strings.ToLower(env))
}

var lookupEnv = os.LookupEnv

// Env
type Env string

// String converts environment value to predefined strings.
func (e Env) String() string {
	return string(e)
}

// IsDev is a helper function to check whether app is running in Dev mode
func (e Env) IsDev() bool {
	return strings.HasPrefix(strings.ToLower(string(e)), envDev)
}
