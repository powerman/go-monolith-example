package cobrax

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/migrate"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/structlog"
	"github.com/spf13/cobra"
)

// gooseUsageTemplate is cobra usage template for goose commands.
const gooseUsageTemplate = `Usage:
  {{.CommandPath}} [command]{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}

Available Commands:
  up                   Migrate the DB to the most recent version available
  up-by-one            Migrate the DB up by 1
  up-to VERSION        Migrate the DB to a specific VERSION
  down                 Roll back the version by 1
  down-to VERSION      Roll back to a specific VERSION
  redo                 Re-run the latest migration
  reset                Roll back all migrations
  status               Dump the migration status for the current DB
  version              Print the current version of the database
  create NAME [sql|go] Creates new migration file with the current timestamp
  fix                  Apply sequential ordering to migrations{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}
`

//nolint:gochecknoglobals // Regexp.
var reGooseCommand = regexp.MustCompile(`^(?:up|up-by-one|up-to\s+\d+|down|down-to\s+\d+|redo|reset|status|version|create\s+\S+\s+(?:go|sql)|fix)$`)

// validGooseCommand returns true if command is a valid goose command.
func validGooseCommand(command string) bool {
	return reGooseCommand.MatchString(command)
}

// GooseConfig contain configuration for goose command.
type GooseConfig struct {
	MySQLConfig *mysql.Config
	GooseDir    string
}

// NewGooseCmd creates new goose command executed by run.
func NewGooseCmd(serviceName string, goose *goosepkg.Instance, getCfg func() (*GooseConfig, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goose",
		Short: "Migrate database schema",
		Args: func(cmd *cobra.Command, args []string) error {
			gooseCmd := strings.Join(args, " ")
			if gooseCmd == "" {
				return ErrRequireFlagsOrCommand
			} else if !validGooseCommand(gooseCmd) {
				return fmt.Errorf("invalid goose command: %s", gooseCmd)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			gooseCmd := strings.Join(args, " ")

			ctx := def.NewContext(serviceName)
			log := structlog.FromContext(ctx, nil)
			cfg, err := getCfg()
			if err != nil {
				log.Fatalf("failed to get config: %s", err)
			}

			err = migrate.Run(ctx, goose, cfg.GooseDir, gooseCmd, cfg.MySQLConfig)
			if err != nil {
				log.Fatalf("failed to run goose %s: %s", gooseCmd, err)
			}
		},
	}
	cmd.SetUsageTemplate(gooseUsageTemplate)
	return cmd
}
