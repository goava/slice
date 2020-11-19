package slice

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
