<%=licenseText%>
package util

import cli "gopkg.in/urfave/cli.v1"

// TraceOptions are the commonly used options when using a Trace service.
type TraceOptions struct {
	Trace         bool
	TraceEndpoint string
}

// ParseTraceOptions fetches the values from urfave/cli Context and
// returns them as a TraceOptions. Uses the names as specified in
// TraceFlags.
func ParseTraceOptions(c *cli.Context) *TraceOptions {
	return &TraceOptions{
		Trace:         c.Bool("trace"),
		TraceEndpoint: c.String("trace-endpoint"),
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
			Usage:  "Endpoint for the trace service",
			EnvVar: "TRACE_ENDPOINT",
		},
	}
}
