<%=licenseText%>
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"<%=repoUrl%>/pkg/util"
	"<%=repoUrl%>/pkg/log"
	"<%=repoUrl%>/pkg/api"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ocgrpc"
	"gopkg.in/urfave/cli.v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var gatewayCommand = cli.Command{
	Name:   "gateway",
	Usage:  "Start gRPC gateway",
	Action: gatewayAction,
	Flags:  append(gatewayFlags, util.TraceFlags()...),
}

var gatewayFlags = []cli.Flag{
	cli.IntFlag{
		Name:   "port",
		Value:  19990,
		EnvVar: "HTTP_PORT",
	},
	cli.StringFlag{
		Name:   "host",
		Value:  "localhost:19991",
		EnvVar: "GRPC_HOST",
	},
}

var gatewayAction = func(c *cli.Context) error {
	log.Info("Starting <%=serviceName%> gateway")

	log.Debug("Parsing gateway options")
	o, err := parseGatewayOptions(c)
	if err != nil {
		log.WithError(err).Error("Unable to validate arguments")
		return errorExitCode
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true})) // grpc-gateway

	// The following handlers will be called in reversed order (ie. bottom to top)
	handler := newHTTPHandler(mux, o.TraceOptions)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	}

	err = api.Register<%=servicePName%>HandlerFromEndpoint(ctx, mux, o.Host, opts)
	if err != nil {
		log.WithError(err).Error("Unable to register handler from Endpoint")
		return errorExitCode
	}

	errc := make(chan error, 2)

	// Shutdown on SIGINT, SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", o.Port),
		Handler: handler,
	}
	
	// Start Gateway server in separate goroutine
	go func() {
		log.WithFields(log.Fields{
			"port":     o.Port,
			"grpcHost": o.Host}).Info("Starting gateway server")
		err := s.ListenAndServe()
		errc <- errors.Wrap(err, "gateway returned an error")
	}()

	err = <-errc
	log.WithError(err).Info("Shutting down")

	// Gracefully shutdown the Gateway server
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Error("An error happened while shutting down")
	}
	return nil
}

type gatewayOptions struct {
	util.TraceOptions

	Port int
	Host string
}

func parseGatewayOptions(c *cli.Context) (*gatewayOptions, error) {
	traceOptions := util.ParseTraceOptions(c)
	if traceOptions.TraceNamespace == "" {
		traceOptions.TraceNamespace = "<%=serviceName%>_gw"
	}

	port := c.Int("port")
	if !validPortNumber(port) {
		return nil, fmt.Errorf("Invalid port number: %d", port)
	}

	return &gatewayOptions{
		TraceOptions: traceOptions,

		Port: port,
		Host: c.String("host"),
	}, nil
}
