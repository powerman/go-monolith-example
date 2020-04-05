// Monolith with embedded microservices.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/powerman/go-monolith-example/internal/concurrent"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/flags"
	"github.com/powerman/go-monolith-example/ms/example"
	"github.com/powerman/structlog"
	"github.com/spf13/cobra"
)

const (
	connectTimeout = 3 * time.Second // must be less than swarm's deploy.update_config.monitor
	shutdownDelay  = 9 * time.Second // `docker stop` use 10s between SIGTERM and SIGKILL
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

type embeddedService interface {
	Name() string
	Init(cmd, serveCmd *cobra.Command)
	Serve(ctxSetup, ctxShutdown Ctx, shutdown func()) error
}

//nolint:gochecknoglobals // Main.
var (
	embeddedServices = []embeddedService{
		Service{},
		example.Service{},
	}

	exe   = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	bi, _ = debug.ReadBuildInfo()
	ver   = bi.Main.Version
	log   = structlog.New()

	logLevel string

	rootCmd = &cobra.Command{
		Use:     exe,
		Short:   "Monolith with embedded microservices",
		Version: ver,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			structlog.DefaultLogger.SetLogLevel(structlog.ParseLevel(logLevel))
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the monolith version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(exe, "version", ver, runtime.Version())
		},
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts embedded microservices",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := runServe()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	msCmd = &cobra.Command{
		Use:   "ms",
		Short: "Run given embedded microservice's command",
		Args:  cobra.NoArgs,
	}
)

func main() {
	err := def.Init()
	if err != nil {
		log.Fatalln("failed to get defaults:", err)
	}

	log.SetDefaultKeyvals(structlog.KeyUnit, "main")

	flags.OneOfStringVar(rootCmd.PersistentFlags(),
		&logLevel, "log.level", []string{"debug", "info", "warn", "err"}, "log level")
	rootCmd.AddCommand(versionCmd, serveCmd, msCmd)

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
			Args:  cobra.NoArgs,
		}
		service.Init(cmd, serveCmd)
		msCmd.AddCommand(cmd)
	}

	_ = rootCmd.Execute()
}

func runServe() error {
	log.Info("started", "version", ver)
	defer log.Info("finished", "version", ver)

	ctxSetup, cancelConnect := context.WithTimeout(context.Background(), connectTimeout)
	defer cancelConnect()

	ctxShutdown, shutdown := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() { <-sigc; shutdown() }()
	go forceShutdown(ctxShutdown)

	services := make([]func(Ctx) error, len(embeddedServices))
	for i := range embeddedServices {
		serve := embeddedServices[i].Serve
		log := structlog.New(structlog.KeyApp, embeddedServices[i].Name())
		ctxSetup := structlog.NewContext(ctxSetup, log) //nolint:govet // Shadow.
		services[i] = func(ctx Ctx) error {
			ctx = structlog.NewContext(ctx, log)
			return serve(ctxSetup, ctx, shutdown)
		}
	}
	return concurrent.Serve(ctxShutdown, shutdown, services...)
}

func forceShutdown(ctx Ctx) {
	<-ctx.Done()
	time.Sleep(shutdownDelay)
	log.Fatalln("failed to graceful shutdown", "version", ver)
}
