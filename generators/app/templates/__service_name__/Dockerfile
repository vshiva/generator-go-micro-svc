FROM golang:1.11 as builder
ARG LD_FLAGS

COPY ./ /go/src/<%=repoUrl%>/
WORKDIR /go/src/<%=repoUrl%>/

RUN go get -u github.com/golang/dep/cmd/dep && dep ensure && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "${LD_FLAGS}" -o bin/<%=pkgName%> ./cmd/<%=serviceName%>

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/<%=repoUrl%>/bin/<%=pkgName%> /usr/bin
ENTRYPOINT ["/usr/bin/<%=pkgName%>"]
CMD ["-h"]
