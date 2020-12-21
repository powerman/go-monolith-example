package def

import (
	"fmt"

	"github.com/powerman/structlog"
	"google.golang.org/grpc/grpclog"
)

// Log field names.
const (
	LogServer   = "server"   // "OpenAPI", "gRPC", "Prometheus metrics", etc.
	LogRemoteIP = "remoteIP" // IP address.
	LogAddr     = "addr"     // host:port.
	LogHost     = "host"     // DNS hostname or IPv4/IPv6 address.
	LogPort     = "port"     // TCP/UDP port number.
	LogFunc     = "func"     // RPC/event handler method name, REST resource path.
	LogUserName = "userName"
	LogGRPCCode = "grpcCode"
)

func setupLog() {
	structlog.DefaultLogger.
		AppendPrefixKeys(
			LogRemoteIP,
			LogGRPCCode,
			LogFunc,
		).
		SetSuffixKeys(
			LogServer,
			LogUserName,
			structlog.KeyStack,
		).
		SetDefaultKeyvals(
			structlog.KeyPID, nil,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyApp:  " %12.12[2]s:", // set to max microservice name length
			structlog.KeyUnit: " %9.9[2]s:",   // set to max KeyUnit/package length
			LogRemoteIP:       " %-15[2]s",    // set to 19.19 or 39 or 45 for IPv6
			LogGRPCCode:       " %-16.16[2]s",
			LogFunc:           " %[2]s:",
			LogHost:           " %[2]s",
			LogPort:           ":%[2]v",
			LogAddr:           " %[2]s",
			"version":         " %s %v",
			"json":            " %s=%#q",
			"ptr":             " %[2]p",   // for debugging references
			"data":            " %#+[2]v", // for debugging structs
			"offset":          " page=%3[2]d",
			"limit":           "+%[2]d ",
			"err":             " %s: %v",
			LogServer:         " [%[2]s]",
			LogUserName:       " %[2]v",
		})
	grpclog.SetLoggerV2(grpcLog{structlog.New(structlog.KeyUnit, "grpcpkg")})
}

type grpcLog struct{ *structlog.Logger }

func (g grpcLog) Info(args ...interface{})                    { g.Debug(fmt.Sprint(args...)) }
func (g grpcLog) Infoln(args ...interface{})                  { g.Debug(g.sprintln(args...)) }
func (g grpcLog) Infof(format string, args ...interface{})    { g.Debug(fmt.Sprintf(format, args...)) }
func (g grpcLog) Warning(args ...interface{})                 { g.Warn(fmt.Sprint(args...)) }
func (g grpcLog) Warningln(args ...interface{})               { g.Warn(g.sprintln(args...)) }
func (g grpcLog) Warningf(format string, args ...interface{}) { g.Warn(fmt.Sprintf(format, args...)) }
func (g grpcLog) Error(args ...interface{})                   { g.PrintErr(fmt.Sprint(args...)) }
func (g grpcLog) Errorln(args ...interface{})                 { g.PrintErr(g.sprintln(args...)) }
func (g grpcLog) Errorf(format string, args ...interface{})   { g.PrintErr(fmt.Sprintf(format, args...)) }
func (g grpcLog) V(l int) bool                                { return false }
func (g grpcLog) sprintln(args ...interface{}) string {
	s := fmt.Sprintln(args...)
	return s[:len(s)-1]
}
