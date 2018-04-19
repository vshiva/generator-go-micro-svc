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
- enables consistent logging

### Building 
#### Pre-Req
- Go 
- Docker

#### Binary
`make`

#### Docker Image
`make docker-build`

### Running

After a successful build using make you

`./bin/<%=pkgName%>`

Similarly you can the docker image after 

`docker run --rm <%=pkgName%>:dev`

### Development workflow
- Add your project specific flags / command line arguments in `cmd/<%=servicePName%>/server.go`
- Define your protobuf messaages and grpc services in `<%=servicePName%>.proto` file.
- Run `make gen`. This will generate requred structs and service methods from the proto file.
- Implement the business logic for your service methods in `pkg/server` package
- Run `make build`. This will build a binary `<%=pkgName%>` under `./bin` folder
- Often times the service needs to interact with some store. For example a sql store or a nosql store like mongo or bolddb. A Store interface defines all the persistence/repo methods and you there could be multiple implementations for the persistence layer. The store interface methods and various implmentations are defined in `pkg/state` package.
 


