// Package cobrax contains helpers to use with github.com/spf13/cobra.
package cobrax

import (
	"errors"

	"github.com/spf13/cobra"
)

// Errors.
var (
	ErrRequireFlagOrCommand = errors.New("require flag or command")
)

// RequireFlagOrCommand should be used as cobra.Command.RunE for "empty"
// commands which are just a containers for subcommands.
func RequireFlagOrCommand(_ *cobra.Command, _ []string) error {
	return ErrRequireFlagOrCommand
}
