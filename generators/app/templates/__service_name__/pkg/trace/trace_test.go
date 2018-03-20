<%=licenseText%>
package trace

import (
	"context"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	zipkintracer "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/types"
	"github.com/stretchr/testify/assert"
)

func Test_ExtractTraceID(t *testing.T) {
	// Only set TraceID, as the rest is ignored
	zipkinSpanContext := zipkintracer.SpanContext{TraceID: types.TraceID{7777, 3333}}

	tests := []struct {
		name     string
		context  context.Context
		expected string
	}{
		{"Nil", nil, ""},
		{"NoSpan", context.Background(), ""},
		{"NoSpanContext", opentracing.ContextWithSpan(context.Background(), &fakeSpan{}), ""},
		{"UnknownSpanContext", opentracing.ContextWithSpan(context.Background(), &fakeSpan{&fakeSpanContext{}}), ""},
		{"ZipkinSpanContext", opentracing.ContextWithSpan(context.Background(), &fakeSpan{zipkinSpanContext}), "1e610000000000000d05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ExtractTraceID(tt.context)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

type fakeSpan struct {
	spanContext opentracing.SpanContext
}

func (f *fakeSpan) Context() opentracing.SpanContext                       { return f.spanContext }
func (f *fakeSpan) SetBaggageItem(key, val string) opentracing.Span        { return nil }
func (f *fakeSpan) BaggageItem(key string) string                          { return "" }
func (f *fakeSpan) SetTag(key string, value interface{}) opentracing.Span  { return f }
func (f *fakeSpan) LogFields(fields ...log.Field)                          {}
func (f *fakeSpan) LogKV(keyVals ...interface{})                           {}
func (f *fakeSpan) Finish()                                                {}
func (f *fakeSpan) FinishWithOptions(opts opentracing.FinishOptions)       {}
func (f *fakeSpan) SetOperationName(operationName string) opentracing.Span { return f }
func (f *fakeSpan) Tracer() opentracing.Tracer                             { return nil }
func (f *fakeSpan) LogEvent(event string)                                  {}
func (f *fakeSpan) LogEventWithPayload(event string, payload interface{})  {}
func (f *fakeSpan) Log(data opentracing.LogData)                           {}

var _ opentracing.Span = (*fakeSpan)(nil)

type fakeSpanContext struct{}

func (f *fakeSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}

var _ opentracing.SpanContext = (*fakeSpanContext)(nil)
