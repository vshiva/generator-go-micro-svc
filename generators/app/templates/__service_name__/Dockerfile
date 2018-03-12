FROM golang:1.10 as builder
WORKDIR /go/src/<%=repoUrl%>/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o <%=serviceName%> .

FROM alpine:3.7  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/<%=repoUrl%>/<%=serviceName%> .
CMD ["./<%=serviceName%>"]  
