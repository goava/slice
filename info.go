package slice

import "strings"

const (
	devPrefix string = "dev"
)

// The application information.
type Info struct {
	// The Application name.
	Name string
	// The raw value of APP_ENV environment variable. Due to the abundance of possible
	// application launch environments, this is a raw value.
	Env Env
	// The debug flag.
	Debug bool
}

// IsDev is a helper function to check whether app is running in Dev mode
func (info *Info) IsDev() bool {
	return strings.HasPrefix(strings.ToLower(string(info.Env)), devPrefix)
}
