package slice

import (
	"os"
	"strings"
)

const (
	// Production environment
	Prod Env = "prod"
	// Development environment
	Dev Env = "dev"
	// Testing environment
	Test Env = "test"
)

func parseEnv(env string) Env {
	switch strings.ToLower(env) {
	case "dev", "develop", "development":
		return Dev
	case "test", "testing":
		return Test
	default:
		return Prod
	}
}

var lookupEnv = os.LookupEnv

// Env
type Env string

// String converts environment value to predefined strings.
func (e Env) String() string {
	return string(e)
}
