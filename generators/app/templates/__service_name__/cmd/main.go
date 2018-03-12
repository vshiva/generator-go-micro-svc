<%=licenseText%>
package main

import (
	"os"

	"github.com/wercker/pkg/log"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "<%=servicePName%>"
	app.Copyright = "(c) 2018 Copyright"
	app.Usage = "<%=servicePName%> description"

	app.Version = Version()
	app.Compiled = CompiledAt()
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
