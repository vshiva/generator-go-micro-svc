<%=licenseText%>
package trace

import (
	"testing"

	"golang.org/x/net/context"

	opentracing "github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"<%=repoUrl%>/pkg/log"
)

func Test_ExposeInterceptor(t *testing.T) {
	// Handler which will extract log Fields (from the context)
	called := false
	var actualFields log.Fields
	th := func(ctx context.Context, req interface{}) (interface{}, error) {
		called = true
		actualFields, _ = log.FieldsFromContext(ctx)

		return nil, nil
	}

	i := ExposeInterceptor()

	// The context and the span context which used for the request
	zipkinSpanContext := zipkintracer.SpanContext{TraceID: types.TraceID{7777, 3333}}
	ctx := opentracing.ContextWithSpan(context.Background(), &fakeSpan{zipkinSpanContext})

	_, err := i(ctx, nil, nil, th)
	require.NoError(t, err)

	// Test that the TraceID was set in the context with the correct value
	if assert.NotNil(t, actualFields) {
		f, ok := actualFields[TraceFieldKey]
		if assert.True(t, ok, "Fields does not contain expected field with key: %s", TraceFieldKey) {
			assert.Equal(t, "1e610000000000000d05", f)
		}
	}
}
