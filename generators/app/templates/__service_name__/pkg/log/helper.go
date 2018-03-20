<%=licenseText%>
package log

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

//SetupLogging helps with setting up logger
func SetupLogging(c *cli.Context) error {
	if c.GlobalBool("debug") {
		SetLevel(DebugLevel)
	}

	// Dynamically return false or true based on the logger output's
	// file descriptor referring to a terminal or not.
	if os.Getenv("TERM") == "dumb" || !isLogrusTerminal() {
		SetFormatter(&logrus.JSONFormatter{})
	}
	return nil
}

// isLogrusTerminal checks if the standard logger of Logrus is a terminal.
func isLogrusTerminal() bool {
	w := logrus.StandardLogger().Out
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
