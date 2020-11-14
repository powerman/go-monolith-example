package cobrax

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/spf13/cobra"

	"github.com/powerman/go-monolith-example/pkg/migrate"
)

// GooseMySQLConfig contain configuration for goose command.
type GooseMySQLConfig struct {
	MySQL         *mysql.Config
	GooseMySQLDir string
}

// NewGooseMySQLCmd creates new goose command executed by run.
func NewGooseMySQLCmd(ctx context.Context, goose *goosepkg.Instance, getCfg func() (*GooseMySQLConfig, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goose-mysql",
		Short: "Migrate MySQL database schema",
		Args:  gooseArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			gooseCmd := strings.Join(args, " ")

			cfg, err := getCfg()
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}

			connector := &migrate.MySQL{Config: cfg.MySQL}
			err = migrate.Run(ctx, goose, cfg.GooseMySQLDir, gooseCmd, connector)
			if err != nil {
				return fmt.Errorf("failed to run goose %s: %w", gooseCmd, err)
			}
			return nil
		},
	}
	cmd.SetUsageTemplate(gooseUsageTemplate)
	return cmd
}
