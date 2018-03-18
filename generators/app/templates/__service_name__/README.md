# <%=servicePName%>

<%=servicePName%> description

## Getting Started
Go Micro service based on a blueprint that 
- enables fast development of gRPC based micro services
- exposes the gRPC services as REST / Json via grpc gateway interface
- exposes metrics endpoint, which prometheus could scrape from
- support open tracing / zipkin
- exposes health check end points
- defiens state management interface. A default implementation of BoltDB
- provides standard CLI 
- build Docker image
- enables Kubernetes friendly deployment using Helm chart
- enables consistent loggin

### Building 
#### Pre-Req
- Go 
- Docker

#### Binary
`make`

#### Docker Image
`make docker-build`

