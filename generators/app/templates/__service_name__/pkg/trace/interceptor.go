<%=licenseText%>
package trace

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	grpcmw "github.com/mwitkow/go-grpc-middleware"
	opentracing "github.com/opentracing/opentracing-go"
	"<%=repoUrl%>/pkg/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Interceptor adds a opentracing middleware, and exposes the TraceID.
func Interceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return grpcmw.ChainUnaryServer(
		otgrpc.OpenTracingServerInterceptor(tracer), // opentracing (incoming)
		ExposeInterceptor(),                         // expose traceID
	)
}

// ExposeInterceptor extracts the TraceID from the context and adds it to
// fields in the context.
func ExposeInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (interface{}, error) {
		traceID := ExtractTraceID(ctx)
		if traceID != "" {
			ctx, _ = log.AddFieldToCtx(ctx, TraceFieldKey, traceID)
		}

		return next(ctx, req)
	}
}
