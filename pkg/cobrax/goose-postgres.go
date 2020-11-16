package cobrax

import (
	"context"
	"fmt"
	"strings"

	goosepkg "github.com/powerman/goose/v2"
	"github.com/spf13/cobra"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/migrate"
)

// GoosePostgresConfig contain configuration for goose command.
type GoosePostgresConfig struct {
	Postgres         *def.PostgresConfig
	GoosePostgresDir string
}

// NewGoosePostgresCmd creates new goose command executed by run.
func NewGoosePostgresCmd(ctx context.Context, goose *goosepkg.Instance, getCfg func() (*GoosePostgresConfig, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goose-postgres",
		Short: "Migrate PostgreSQL database schema",
		Args:  gooseArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			gooseCmd := strings.Join(args, " ")

			cfg, err := getCfg()
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}

			connector := &migrate.Postgres{PostgresConfig: cfg.Postgres}
			err = migrate.Run(ctx, goose, cfg.GoosePostgresDir, gooseCmd, connector)
			if err != nil {
				return fmt.Errorf("failed to run goose %s: %w", gooseCmd, err)
			}
			return nil
		},
	}
	cmd.SetUsageTemplate(gooseUsageTemplate)
	return cmd
}
