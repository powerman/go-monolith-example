// Package cobrax contains helpers to use with github.com/spf13/cobra.
package cobrax

import (
	"errors"

	"github.com/spf13/cobra"
)

// Errors.
var (
	ErrRequireFlagsOrCommand = errors.New("require flags or command")
)

// RequireFlagsOrCommand should be used as cobra.Command.RunE for "empty"
// commands which are just a containers for subcommands.
func RequireFlagsOrCommand(cmd *cobra.Command, args []string) error {
	return ErrRequireFlagsOrCommand
}
