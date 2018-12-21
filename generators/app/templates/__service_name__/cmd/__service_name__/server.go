<%=licenseText%>
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"<%=repoUrl%>/pkg/api"
	"<%=repoUrl%>/pkg/server"
	"<%=repoUrl%>/pkg/state"
	"<%=repoUrl%>/pkg/util"
	"<%=repoUrl%>/pkg/health"
	"<%=repoUrl%>/pkg/log"
	cli "gopkg.in/urfave/cli.v1"
)

var serverCommand = cli.Command{
	Name:   "server",
	Usage:  "start gRPC server",
	Action: serverAction,
	Flags:  append(serverFlags, util.TraceFlags()...),
}

var serverFlags = []cli.Flag{
	cli.IntFlag{
		Name:   "port",
		Value:  19991,
		EnvVar: "PORT",
	},
	cli.IntFlag{
		Name:   "health-port",
		Value:  19992,
		EnvVar: "HEALTH_PORT",
	},
	cli.IntFlag{
		Name:   "metrics-port",
		Value:  9102,
		EnvVar: "METRICS_PORT",
	},
	cli.StringFlag{
		Name:  "state-store",
		Usage: "storage driver, currently supported [memory]",
		Value: "memory",
	},
}

var serverAction = func(c *cli.Context) error {
	log.Info("Starting <%=serviceName%> server")

	log.Debug("Parsing server options")
	o, err := parseServerOptions(c)
	if err != nil {
		log.WithError(err).Error("Unable to validate arguments")
		return errorExitCode
	}

	mux := http.NewServeMux()
	registerPromMetricsExporter(mux, "<%=serviceName%>_server")
	registerTraceExporter(o.TraceOptions)

	healthService := health.New()

	store, err := getStore(o)
	if err != nil {
		log.WithError(err).Error("Unable to create store")
		return errorExitCode
	}
	defer store.Close()

	err = store.Initialize(context.Background())
	if err != nil {
		log.WithError(err).Error("Unable to initialize store")
		return errorExitCode
	}

	healthService.RegisterProbe("store", store)

	store = state.NewTraceStore(store)   // Wrap it with tracing
	store = state.NewMetricsStore(store) // Wrap it with metrics.

	log.Debug("Creating server")
	server, err := server.New(store)
	if err != nil {
		log.WithError(err).Error("Unable to create server")
		return errorExitCode
	}

	s := newGRPCServeWithMetrics()
	api.Register<%=servicePName%>Server(s, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", o.Port))
	if err != nil {
		log.WithField("port", o.Port).WithError(err).Error("Failed to listen")
		return errorExitCode
	}

	errc := make(chan error, 4)

	// Shutdown on SIGINT, SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Start gRPC server
	go func() {
		log.WithField("port", o.Port).Info("Starting server")
		err := s.Serve(lis)
		errc <- errors.Wrap(err, "server returned an error")
	}()

	// Start health server
	go func() {
		log.WithField("port", o.HealthPort).Info("Starting health service")
		err := healthService.ListenAndServe(fmt.Sprintf(":%d", o.HealthPort))
		errc <- errors.Wrap(err, "health service returned an error")
	}()

	// Start metrics server
	go func() {
		log.WithField("port", o.MetricsPort).Info("Starting metrics server")
		err := http.ListenAndServe(fmt.Sprintf(":%d", o.MetricsPort), mux)
		errc <- errors.Wrap(err, "metrics service returned an error")
	}()

	err = <-errc
	log.WithError(err).Info("Shutting down")

	// Gracefully shutdown the health server
	healthService.Shutdown(context.Background())

	// Gracefully shutdown the gRPC server
	s.GracefulStop()

	return nil
}

type serverOptions struct {
	util.TraceOptions

	Port          int
	HealthPort    int
	MetricsPort   int
	StateStore    string
}

func parseServerOptions(c *cli.Context) (*serverOptions, error) {
	traceOptions := util.ParseTraceOptions(c)
	if traceOptions.TraceNamespace == "" {
		traceOptions.TraceNamespace = "<%=serviceName%>_server"
	}

	port := c.Int("port")
	if !validPortNumber(port) {
		return nil, fmt.Errorf("invalid port number: %d", port)
	}

	healthPort := c.Int("health-port")
	if !validPortNumber(healthPort) {
		return nil, fmt.Errorf("invalid health-port number: %d", healthPort)
	}

	if healthPort == port {
		return nil, errors.New("health-port and port cannot be the same")
	}

	metricsPort := c.Int("metrics-port")
	if !validPortNumber(metricsPort) {
		return nil, fmt.Errorf("invalid metrics port number: %d", metricsPort)
	}

	if metricsPort == port {
		return nil, errors.New("metrics-port and port cannot be the same")
	}

	if metricsPort == healthPort {
		return nil, errors.New("metrics-port and health-port cannot be the same")
	}

	return &serverOptions{
		TraceOptions: traceOptions,
		
		Port:          port,
		HealthPort:    healthPort,
		MetricsPort:   metricsPort,
		StateStore:    c.String("state-store"),
	}, nil
}

func getStore(o *serverOptions) (state.Store, error) {
	switch o.StateStore {
	case "memory":
		return state.NewInMemoryStore(), nil
	default:
		return nil, fmt.Errorf("Invalid store: %s", o.StateStore)
	}
}
