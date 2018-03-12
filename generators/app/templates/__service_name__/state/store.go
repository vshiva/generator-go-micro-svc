<%=licenseText%>
package state

import (
	"io"

	"github.com/wercker/pkg/health"
)

// Store provides access to data that is required for <%=licenseText%>.
type Store interface {
	Initialize() error
	io.Closer
	health.Probe
}
