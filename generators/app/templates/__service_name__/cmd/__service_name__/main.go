<%=licenseText%>
package main

import (
	"os"
	"fmt"
	"runtime"
	"<%=repoUrl%>/pkg/log"
	cli "gopkg.in/urfave/cli.v1"
	"github.com/opentracing/opentracing-go"
	"<%=repoUrl%>/pkg/util"
	svc "<%=repoUrl%>"
	"<%=repoUrl%>/pkg/trace"
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

func getTracer(o *util.TraceOptions, serviceName string, port int) (opentracing.Tracer, error) {
	if o.Trace {
		return trace.NewZipkinTracer(o.TraceEndpoint, serviceName, port)
	}

	return trace.NewNoopTracer()
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
