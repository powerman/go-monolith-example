package def

import (
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/powerman/getenv"
)

// Default ports to serve metrics.
const (
	MetricsPort = 16000 + iota
	ExampleMetricsPort
)

// Default ports to serve RPC.
const (
	_ = 17000 + iota
	ExampleRPCPort
)

// Constants.
var (
	ProgName              = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	Hostname, hostnameErr = os.Hostname()
	testTimeFactor        = getenv.Float("GO_TEST_TIME_FACTOR", 1.0)
	TestSecond            = time.Duration(float64(time.Second) * testTimeFactor)
)

// Version returns application version based on build info.
func Version() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Main.Version
	}
	return "(test)"
}
