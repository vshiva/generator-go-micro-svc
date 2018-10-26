<%=licenseText%>
package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"<%=repoUrl%>/pkg/util"
	"<%=repoUrl%>/pkg/log"
	svc "<%=repoUrl%>"
	
	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	// errorExitCode returns a urfave decorated error which indicates a exit
	// code 1. To be returned from a urfave action.
	errorExitCode = cli.NewExitError("", 1)
)

// validPortNumber returns true if port is between 0 and 65535.
func validPortNumber(port int) bool {
	return port > 0 && port < 65535
}

// registerTraceExporter registers trace exportre
func registerTraceExporter(options util.TraceOptions) {

	if !options.Trace {
		log.Debugf("Tracing not enabled")
		return
	}

	var exporter trace.Exporter
	var err error
	switch options.TraceBackend {
	case "jaeger":
		// Register the Jaeger exporter to be able to retrieve
		// the collected spans.
		exporter, err = jaeger.NewExporter(jaeger.Options{
			CollectorEndpoint: options.TraceEndpoint,
			Process: jaeger.Process{
				ServiceName: options.TraceNamespace,
			},
		})
	default:
		log.Warnf("Unsupported tracing backend %s", options.TraceBackend)
		return
	}

	if err != nil {
		log.Fatalf("Failed to create an Jaeger Trace exporter %v", err)
		return
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// grpc Server Tracers
	// http Server Middleware Tracers

}

func registerPromMetricsExporter(mux *http.ServeMux, serviceName string) {

	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatalf("Failed to register ocgrpc server views: %v", err)
	}

	// Create the Prometheus exporter.
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: serviceName,
	})
	if err != nil {
		log.Fatalf("Failed to create prometheus metrics exporter: %v", err)
	}

	view.RegisterExporter(pe)
	log.Debug("Registering prometheus exporter with http server mux")

	mux.Handle("/metrics", pe)
	zpages.Handle(mux, "/")

}

func newGRPCServeWithMetrics() *grpc.Server {
	log.Debug("Creating new gRPC server with default server metrics views")

	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatalf("Failed to register ocgrpc server views: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))

	return s
}

func newHTTPHandler(mux *grpc_runtime.ServeMux, options util.TraceOptions) http.Handler {
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatalf("Failed to register ochttp server views: %v", err)
	}
	// Register views to collect data.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}
	registerTraceExporter(options)
	return &ochttp.Handler{
		Handler: mux,
	}
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n Version:  %s\n Git Commit:  %s\n Go Version:  %s\n OS/Arch:  %s/%s\n Built:  %s\n",
			c.App.Name, c.App.Version, svc.GitCommit,
			runtime.Version(), runtime.GOOS, runtime.GOARCH,c.App.Compiled.String())
	}

	app := cli.NewApp()

	app.Name = "<%=pkgName%>"
	app.Copyright = "(c) 2018 Copyright"
	app.Usage = "<%=servicePName%> description"

	app.Version = svc.Version()
	app.Compiled = svc.CompiledAt()
	app.Before = log.SetupLogging
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
	}
	app.Commands = []cli.Command{
		gatewayCommand,
		serverCommand,
	}

	app.Run(os.Args)
}
