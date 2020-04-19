package config

import "github.com/spf13/pflag"

// AddFlag defines a flag with the specified name and usage string.
// Calling it again with same fs, value and name will have no effect.
func AddFlag(fs *pflag.FlagSet, value pflag.Value, name string, usage string) {
	if flag := fs.Lookup(name); flag != nil && flag.Value == value {
		return
	}
	fs.Var(value, name, usage)
}
