<%=licenseText%>
package state

import (
	"io"
	"context"
	"<%=repoUrl%>/pkg/health"
)

// Store provides access to data that is required for <%=licenseText%>.
type Store interface {
	Initialize(ctx context.Context) error
	io.Closer
	health.Probe
}
