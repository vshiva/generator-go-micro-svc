#!/bin/bash
set -e

(

LOCAL=$(dirname $PWD)
ROOT=$LOCAL
protoc="protoc \
  -I/usr/local/include
  -I.
  -I$GOPATH/src \
  -I./vendor \
  -I./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"

cd $LOCAL

echo "Generating gRPC server, gateway, swagger, flow"
$protoc --go_out=plugins=grpc:$ROOT/<%=serviceName%>pb \
        --grpc-gateway_out=logtostderr=true,request_context=true:$ROOT/<%=serviceName%>pb \
        --swagger_out=logtostderr=true:$ROOT/<%=serviceName%>pb \
        --flow_out=$ROOT/<%=serviceName%>pb \
        <%=serviceName%>.proto

)
