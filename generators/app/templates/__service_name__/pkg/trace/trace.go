<%=licenseText%>
package trace

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
)

// ExtractTraceID extracts the TraceID from a opentracing enabled context.
// Currently only the zipkin implementation is supported. Returns an empty
// string when no opentracing span was found, or a unsupported implementation
// was used.
func ExtractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	spancontext := span.Context()
	if spancontext == nil {
		return ""
	}

	switch s := spancontext.(type) {
	case zipkintracer.SpanContext:
		return s.TraceID.ToHex()
	default:
		return ""
	}
}
