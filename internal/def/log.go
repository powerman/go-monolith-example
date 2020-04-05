package def

import (
	"context"

	"github.com/powerman/structlog"
)

// Log field names.
const (
	LogHost    = "host"
	LogPort    = "port"
	LogAddr    = "addr"
	LogRemote  = "remote" // aligned IPv4:Port "   192.168.0.42:1234 "
	LogFunc    = "func"   // RPC method name, REST resource path
	LogUser    = "userID"
	LogService = "service"
)

func setupLog() {
	structlog.DefaultLogger.
		AppendPrefixKeys(
			LogRemote,
			LogFunc,
		).
		SetSuffixKeys(
			LogService,
			LogUser,
			structlog.KeyStack,
		).
		SetDefaultKeyvals(
			structlog.KeyPID, nil,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyApp:  " %12.12[2]s:", // set to max microservice name length
			structlog.KeyUnit: " %7.7[2]s:",   // set to max KeyUnit/package length
			LogHost:           " %[2]s",
			LogPort:           ":%[2]v",
			LogAddr:           " %[2]s",
			LogRemote:         " %-21[2]s",
			LogFunc:           " %[2]s:",
			LogService:        " [%[2]s]",
			LogUser:           " u:%[2]d",
			"version":         " %s %v",
			"err":             " %s: %v",
			"json":            " %s=%#q",
			"offset":          " page=%3[2]d",
			"limit":           "+%[2]d ",
			"ptr":             " %[2]p",   // for debugging references
			"data":            " %#+[2]v", // for debugging structs
		})
}

// NewContext returns context.Background() which contains logger
// configured for given service.
func NewContext(service string) context.Context {
	return structlog.NewContext(context.Background(), structlog.New(structlog.KeyApp, service))
}
