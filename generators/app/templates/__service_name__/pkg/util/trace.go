<%=licenseText%>
package util

import cli "gopkg.in/urfave/cli.v1"

// TraceOptions are the commonly used options when using a Trace service.
type TraceOptions struct {
	Trace          bool
	TraceEndpoint  string
	TraceNamespace string
	TraceBackend   string
}

// ParseTraceOptions fetches the values from urfave/cli Context and
// returns them as a TraceOptions. Uses the names as specified in
// TraceFlags.
func ParseTraceOptions(c *cli.Context) TraceOptions {
	return TraceOptions{
		Trace:          c.Bool("trace"),
		TraceEndpoint:  c.String("trace-endpoint"),
		TraceNamespace: c.String("trace-namespace"),
		TraceBackend:   c.String("trace-backend"),
	}
}

// TraceFlags returns the flags that will be used by
// ParseTraceOptions.
func TraceFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   "trace",
			Usage:  "Enable tracing",
			EnvVar: "TRACE_ENABLED",
		},
		cli.StringFlag{
			Name:   "trace-endpoint",
			Value: "http://localhost:14268/api/traces",
			Usage:  "Endpoint for the trace service",
			EnvVar: "TRACE_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "trace-namespace",
			Usage:  "Servie namespace",
			EnvVar: "TRACE_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "trace-backend",
			Usage:  "Backend to use for the tracing",
			Value:  "jaeger",
			EnvVar: "TRACE_BACKEND",
		},
	}
}
