// Monolith with embedded microservices.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/powerman/appcfg"
	"github.com/powerman/go-monolith-example/internal/cobrax"
	"github.com/powerman/go-monolith-example/internal/concurrent"
	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/ms/example"
	"github.com/powerman/go-monolith-example/ms/metrics"
	"github.com/powerman/structlog"
	"github.com/spf13/cobra"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

type embeddedService interface {
	Name() string
	Init(cfg *config.Cfg, cmd, serveCmd *cobra.Command) error
	Serve(ctxStartup, ctxShutdown Ctx, shutdown func()) error
}

//nolint:gochecknoglobals // Main.
var (
	embeddedServices = []embeddedService{
		metrics.Service{},
		example.Service{},
	}

	log = structlog.New()

	logLevel             = appcfg.MustOneOfString("debug", []string{"debug", "info", "warn", "err"})
	serveStartupTimeout  = appcfg.MustDuration("3s") // must be less than swarm's deploy.update_config.monitor
	serveShutdownTimeout = appcfg.MustDuration("9s") // `docker stop` use 10s between SIGTERM and SIGKILL

	rootCmd = &cobra.Command{
		Use:     def.ProgName,
		Short:   "Monolith with embedded microservices",
		Version: fmt.Sprintf("%s %s", def.Version(), runtime.Version()),
		RunE:    cobrax.RequireFlagsOrCommand,
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts embedded microservices",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if err := runServe(); err != nil {
				log.Fatal(err)
			}
		},
	}
	msCmd = &cobra.Command{
		Use:   "ms",
		Short: "Run given embedded microservice's command",
		RunE:  cobrax.RequireFlagsOrCommand,
	}
)

func main() {
	err := def.Init()
	if err != nil {
		log.Fatalln("failed to get defaults:", err)
	}

	log.SetDefaultKeyvals(structlog.KeyUnit, "main")

	cfg, err := config.Get()
	if err != nil {
		log.Fatalln("failed to init config:", err)
	}

	rootCmd.PersistentFlags().Var(&logLevel, "log.level", "log level [debug|info|warn|err]")
	rootCmd.AddCommand(serveCmd, msCmd)
	cobra.OnInitialize(func() {
		structlog.DefaultLogger.SetLogLevel(structlog.ParseLevel(logLevel.String()))
	})

	serveCmd.Flags().Var(&serveStartupTimeout, "timeout.startup", "must be less than swarm's deploy.update_config.monitor")
	serveCmd.Flags().Var(&serveShutdownTimeout, "timeout.shutdown", "must be less than 10s used by 'docker stop' between SIGTERM and SIGKILL")

	seen := make(map[string]bool)
	for _, service := range embeddedServices {
		name := service.Name()
		if seen[name] {
			panic(fmt.Sprintf("duplicate service: %s", name))
		}
		seen[name] = true

		cmd := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Run %s microservice's command", name),
			RunE:  cobrax.RequireFlagsOrCommand,
		}
		err := service.Init(cfg, cmd, serveCmd)
		if err != nil {
			log.Fatalf("failed to init service %s: %s", name, err)
		}
		msCmd.AddCommand(cmd)
	}

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func runServe() error {
	log.Info("started", "version", def.Version())
	defer log.Info("finished", "version", def.Version())

	ctxStartup, cancel := context.WithTimeout(context.Background(), serveStartupTimeout.Value(nil))
	defer cancel()

	ctxShutdown, shutdown := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() { <-sigc; shutdown() }()
	go forceShutdown(ctxShutdown)

	services := make([]func(Ctx) error, len(embeddedServices))
	for i := range embeddedServices {
		serve := embeddedServices[i].Serve
		log := structlog.New(structlog.KeyApp, embeddedServices[i].Name())
		ctxStartup := structlog.NewContext(ctxStartup, log) //nolint:govet // Shadow.
		services[i] = func(ctxShutdown Ctx) error {
			ctxShutdown = structlog.NewContext(ctxShutdown, log)
			return serve(ctxStartup, ctxShutdown, shutdown)
		}
	}
	return concurrent.Serve(ctxShutdown, shutdown, services...)
}

func forceShutdown(ctxShutdown Ctx) {
	<-ctxShutdown.Done()
	time.Sleep(serveShutdownTimeout.Value(nil))
	log.Fatalln("failed to graceful shutdown", "version", def.Version())
}
