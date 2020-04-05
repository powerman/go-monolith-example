// Package flags provides helpers to use with github.com/spf13/cobra.
package flags

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type pValue struct{ flag.Value }

func (pValue) Type() string { return "" }

// Var works like flags.Var except it allows to set default value.
func Var(cmd *cobra.Command, v flag.Value, name string, value interface{}, usage string) {
	if cmd.Flags().Lookup(name) != nil {
		return
	}

	pv, ok := v.(pflag.Value)
	if !ok {
		pv = pValue{v}
	}
	err := pv.Set(fmt.Sprint(value))
	cmd.Flags().Var(pv, name, usage)
	if err != nil {
		if err := cmd.MarkFlagRequired(name); err != nil {
			panic(err)
		}
	}
}

type oneOfString struct {
	oneOf []string
	value *string
}

func (v *oneOfString) Type() string   { return "string" }
func (v *oneOfString) String() string { return *v.value }
func (v *oneOfString) Set(s string) error {
	for _, item := range v.oneOf {
		if s == item {
			*v.value = s
			return nil
		}
	}
	return fmt.Errorf("must be one of %q", v.oneOf)
}

// OneOfStringVar defines a string flag with specified name, list of
// allowed values (first one will become default value), and usage string.
// The argument p points to a string variable in which to store the value
// of the flag.
func OneOfStringVar(f *pflag.FlagSet, p *string, name string, values []string, usage string) {
	if len(values) == 0 {
		panic("values must not be empty")
	}
	pv := &oneOfString{
		oneOf: values,
		value: p,
	}
	err := pv.Set(values[0])
	if err != nil {
		panic("oneOfString: failed to set default to first of allowed values")
	}
	f.Var(pv, name, fmt.Sprintf("%s %q", usage, values))
}

// NotEmptyString is a non-empty string.
type NotEmptyString string

// Type implements pflags.Value interface.
func (v *NotEmptyString) Type() string { return "string" }

// String implements flags.Value interface.
func (v *NotEmptyString) String() string { return string(*v) }

// Set implements flags.Value interface.
func (v *NotEmptyString) Set(s string) error {
	if s == "" {
		return errors.New("required")
	}
	*v = NotEmptyString(s)
	return nil
}

// Port is an integer between 1 and 65535.
type Port int

// Type implements pflags.Value interface.
func (v *Port) Type() string { return "port" }

// String implements flags.Value interface.
func (v *Port) String() string { return strconv.Itoa(int(*v)) }

// Set implements flags.Value interface.
func (v *Port) Set(s string) error {
	const maxPort = 65536
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	} else if !(0 < i && i < maxPort) {
		return fmt.Errorf("must be between 1 and %d", maxPort-1)
	}
	*v = Port(i)
	return nil
}

// Endpoint is an url with host and without trailing slashes.
type Endpoint string

// Type implements pflags.Value interface.
func (v *Endpoint) Type() string { return "endpoint" }

// String implements flags.Value interface.
func (v *Endpoint) String() string { return string(*v) }

// Set implements flags.Value interface.
func (v *Endpoint) Set(s string) error {
	s = strings.TrimRight(s, "/")
	p, err := url.Parse(s)
	if err != nil {
		return err
	} else if p.Host == "" {
		return errors.New("must contain host")
	}
	*v = Endpoint(s)
	return nil
}
