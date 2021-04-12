package slice

import (
	"os"
	"text/tabwriter"

	"github.com/kelseyhightower/envconfig"
)

// DefaultTableFormat constant to use to display usage in a tabular format
const DefaultTableFormat = `{{range .}}{{usage_key .}}	{{usage_type .}}	{{usage_default .}}	{{usage_required .}}	{{usage_description .}}
{{end}}`

// Parameter contains application configuration data.
//
//	type Parameters struct {
//		Addr         string        `envconfig:"addr"`
//		ReadTimeout  time.Duration `envconfig:"read_timeout"`
//		WriteTimeout time.Duration `envconfig:"write_timeout"`
//	}
type Parameter interface {
}

// ParameterParser parses application parameters.
type ParameterParser interface {
	// Parse parses Parameter fields.
	Parse(prefix string, parameters ...Parameter) error
	// Usage prints parameters usage to stdout.
	Usage(prefix string, parameters ...Parameter) error
}

type stdParameterParser struct {
}

func (d stdParameterParser) Parse(prefix string, parameters ...Parameter) error {
	for _, parameter := range parameters {
		if err := envconfig.Process(prefix, parameter); err != nil {
			return err
		}
	}
	return nil
}

func (d stdParameterParser) Usage(prefix string, parameters ...Parameter) error {
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	if _, err := tabs.Write([]byte("KEY\tTYPE\tDEFAULT\tREQUIRED\tDESCRIPTION\n")); err != nil {
		return err
	}
	for _, parameter := range parameters {
		if err := envconfig.Usagef(prefix, parameter, tabs, DefaultTableFormat); err != nil {
			return err
		}
	}
	return tabs.Flush()
}
