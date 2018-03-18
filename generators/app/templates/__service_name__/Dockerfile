FROM golang:1.10 as builder
ARG LD_FLAGS

COPY ./ /go/src/<%=repoUrl%>/
WORKDIR /go/src/<%=repoUrl%>/

RUN go get -u github.com/kardianos/govendor && govendor sync && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "${LD_FLAGS}" -o <%=serviceName%> ./cmd/<%=serviceName%>

FROM alpine:3.7  
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/<%=repoUrl%>/<%=serviceName%> /usr/bin
ENTRYPOINT ["/usr/bin/<%=serviceName%>"]
CMD ["-h"]
