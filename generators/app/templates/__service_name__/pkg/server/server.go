<%=licenseText%>
package server

import (
	"<%=repoUrl%>/pkg/api"
	"<%=repoUrl%>/pkg/state"

	"golang.org/x/net/context"
)

// New Creates a new <%=servicePName%>Server which implements <%=serviceName%>pb.<%=servicePName%>Server.
func New(store state.Store) (*<%=servicePName%>Server, error) {
	return &<%=servicePName%>Server{
		store: store,
	}, nil
}

// <%=servicePName%>Server implements <%=serviceName%>pb.<%=servicePName%>Server.
type <%=servicePName%>Server struct {
	store state.Store
}

// Action is a example implementation and should be replaced with an actual
// implementation.
func (s *<%=servicePName%>Server) Action(ctx context.Context, req *api.ActionRequest) (*api.ActionResponse, error) {
	return &api.ActionResponse{}, nil
}

// Make sure that <%=servicePName%>Server implements the <%=serviceName%>pb.<%=servicePName%>Service interface.
var _ api.<%=servicePName%>Server = &<%=servicePName%>Server{}
