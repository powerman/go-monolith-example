package cobrax

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

var (
	reGooseCommand    = regexp.MustCompile(`^(?:up|up-by-one|up-to\s+\d+|down|down-to\s+\d+|redo|reset|status|version|create\s+\S+\s+(?:go|sql)|fix)$`) //nolint:gochecknoglobals // Regexp.
	errInvalidCommand = errors.New("invalid goose command")
)

// validGooseCommand returns true if command is a valid goose command.
func validGooseCommand(command string) bool {
	return reGooseCommand.MatchString(command)
}

func gooseArgs(cmd *cobra.Command, args []string) error {
	gooseCmd := strings.Join(args, " ")
	if gooseCmd == "" {
		return ErrRequireFlagOrCommand
	} else if !validGooseCommand(gooseCmd) {
		return fmt.Errorf("%w: %s", errInvalidCommand, gooseCmd)
	}
	return nil
}
