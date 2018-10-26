# <%=servicePName%>

<%=servicePName%> description

A Micro service in GoLang based that
- enables fast development of [gRPC](https://grpc.io/) based micro services
- exposes the gRPC services as REST / Json via [grpc gateway](https://github.com/grpc-ecosystem/grpc-gateway) interface
- exposes metrics endpoint, which [Prometheus](https://prometheus.io/) could scrape from
- support tracing and metrics instrumentation using [OpenCensus](https://opencensus.io/)
- exposes health check end points
- defines state management interface. A default setup for [PostgreSQL](https://www.postgresql.org/)
- provides standard CLI 
- build [Docker](https://www.docker.com/) container image
- enables [Kubernetes](https://kubernetes.io/) friendly deployment using [Helm charts](https://helm.sh/)
- enables consistent logging

## Getting Started

### Building 
#### Pre-Req
- [Go 1.11](https://golang.org/dl/)
- [Docker](https://store.docker.com/search?q=&type=edition&offering=community)
- [Protocol Buffers 3.6.1](https://github.com/protocolbuffers/protobuf/releases)
- *[Kubernetes](https://docs.docker.com/docker-for-mac/kubernetes/)* - `Optional` If you want to deploy

#### Binary
`make`

#### Docker Image
`make docker-build`

### Development workflow
- Add your project specific flags / command line arguments in `cmd/<%=servicePName%>/server.go`
- Define your protobuf messaages and grpc services in `<%=servicePName%>.proto` file.
- Run `make gen`. This will generate requred structs and service methods from the proto file.
- Implement the business logic for your service methods in `pkg/server` package
- Run `make build`. This will build a binary `<%=pkgName%>` under `./bin` folder
- Often times the service needs to interact with some store. For example a sql store or a nosql store like mongo or bolddb. A Store interface defines all the persistence/repo methods and you there could be multiple implementations for the persistence layer. The store interface methods and various implmentations are defined in `pkg/state` package.

### Running

#### on localhost
After a successful build using make you can

`./bin/<%=pkgName%>`


#### on localhost using docker 
Similarly you can the docker image after 

`docker run --rm <%=pkgName%>:dev`

#### on a kubernetes cluster

`todo`