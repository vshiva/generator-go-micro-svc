#!/bin/bash
set -e

(

LOCAL=$(dirname $PWD)/..
ROOT=$LOCAL
protoc="protoc \
    -I. \
    -I$GOPATH/src \
    -I/usr/local/include \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"

cd $LOCAL

echo "Generating gRPC server, gateway, swagger"
$protoc --go_out=plugins=grpc:$ROOT/pkg/api \
        --grpc-gateway_out=logtostderr=true,request_context=true:$ROOT/pkg/api \
        --swagger_out=logtostderr=true:$ROOT/pkg/api \
        <%=serviceName%>.proto

)