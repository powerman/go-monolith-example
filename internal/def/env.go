package def

import (
	"os"
	"time"

	"github.com/powerman/getenv"
)

const (
	mainMetricsPort = 16000 + iota
	exampleMetricsPort
)
const (
	_ = 17000 + iota
	exampleRPCPort
)

// Default values.
var (
	hostname, hostnameErr = os.Hostname()
	TestTimeFactor        = getenv.Float("GO_TEST_TIME_FACTOR", 1.0)
	TestSecond            = time.Duration(float64(time.Second) * TestTimeFactor)
	MetricsHost           = getenv.Str("MONO_METRICS_HOST", hostname)
	MetricsPort           = getenv.Int("MONO_METRICS_PORT", mainMetricsPort)
	RPCHost               = getenv.Str("MONO_RPC_HOST", hostname)
	MySQLHost             = getenv.Str("MONO_MYSQL_HOST", "localhost")
	MySQLPort             = getenv.Int("MONO_MYSQL_PORT", 3306)
	NATSUrls              = getenv.Str("MONO_NATS_URLS", "")
	STANClusterID         = getenv.Str("MONO_STAN_CLUSTER_ID", "")
	ExampleMetricsPort    = getenv.Int("MONO_EXAMPLE_METRICS_PORT", exampleMetricsPort)
	ExampleRPCPort        = getenv.Int("MONO_EXAMPLE_RPC_PORT", exampleRPCPort)
	ExampleDBUser         = getenv.Str("MONO_EXAMPLE_DB_USER", "")
	ExampleDBPass         = getenv.Str("MONO_EXAMPLE_DB_PASS", "")
	ExampleDBName         = getenv.Str("MONO_EXAMPLE_DB_NAME", "example")
	ExampleGooseDir       = "ms/example/internal/migrations"
)
