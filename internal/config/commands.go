package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/powerman/go-monolith-example/internal/flags"
	"github.com/powerman/structlog"
	"github.com/spf13/cobra"
)

// NewGooseCmd creates new goose command executed by run.
func NewGooseCmd(run func(gooseCmd string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goose",
		Short: "Migrate database schema",
		Args: func(cmd *cobra.Command, args []string) error {
			gooseCmd := strings.Join(args, " ")
			if gooseCmd == "" {
				return errors.New("require command")
			} else if !flags.ValidGooseCommand(gooseCmd) {
				return fmt.Errorf("invalid goose command: %s", gooseCmd)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := run(strings.Join(args, " "))
			if err != nil {
				structlog.New().Fatal(err)
			}
		},
	}
	cmd.SetUsageTemplate(flags.GooseUsageTemplate)
	return cmd
}
