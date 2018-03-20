<%=licenseText%>
package trace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"<%=repoUrl%>/pkg/log"
)

func Test_ExposeHandler(t *testing.T) {
	// Handler which will extract log Fields (from the context)
	called := false
	var actualFields log.Fields
	th := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		called = true
		actualFields, _ = log.FieldsFromContext(req.Context())
	})

	// The context and the span context which used for the request
	zipkinSpanContext := zipkintracer.SpanContext{TraceID: types.TraceID{7777, 3333}}
	ctx := opentracing.ContextWithSpan(context.Background(), &fakeSpan{zipkinSpanContext})

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "localhost", nil)
	require.NoError(t, err, "Unable to create request")

	req = req.WithContext(ctx)

	h := ExposeHandler(th)
	h.ServeHTTP(recorder, req)

	require.True(t, called, "Wrapped handler was not called")

	// Test that the TraceID was set in the context with the correct value
	if assert.NotNil(t, actualFields) {
		f, ok := actualFields[TraceFieldKey]
		if assert.True(t, ok, "Fields does not contain expected field with key: %s", TraceFieldKey) {
			assert.Equal(t, "1e610000000000000d05", f)
		}
	}

	// Test that the TraceID was set in response headers with the correct value
	assert.Equal(t, "1e610000000000000d05", recorder.Header().Get(TraceHTTPHeader))
}
