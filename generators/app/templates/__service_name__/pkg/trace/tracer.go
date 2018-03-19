<%=licenseText%>
package trace

import (
	"fmt"
	"os"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
	"github.com/pkg/errors"
)

// NewZipkinTracer creates a new Tracer which uses Zipkin as the backend store.
func NewZipkinTracer(endpoint string, serviceName string, servicePort int) (opentracing.Tracer, error) {
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}

	if strings.Count(endpoint, ":") == 1 {
		endpoint = fmt.Sprintf("%s:9411", endpoint)
	}

	collector, err := zipkintracer.NewHTTPCollector(fmt.Sprintf("%s/api/v1/spans", endpoint))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Zipkin HTTP collector")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "¯\\_(ツ)_/¯"
	}

	debug := false
	recorder := zipkintracer.NewRecorder(collector, debug, fmt.Sprintf("%s:%d", hostname, servicePort), serviceName)

	sameSpan := false
	tracer, err := zipkintracer.NewTracer(recorder, zipkintracer.ClientServerSameSpan(sameSpan))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Zipkin tracer")
	}

	return tracer, nil
}

// NewNoopTracer creates a Tracer which still uses the zipkin Tracer but none
// of the traces will be sampled. This still allows for a unique TraceID to be
// generated.
func NewNoopTracer() (opentracing.Tracer, error) {
	sameSpan := false
	tracer, err := zipkintracer.NewTracer(&noopRecorder{},
		zipkintracer.ClientServerSameSpan(sameSpan),
		zipkintracer.WithSampler(neverSample),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Zipkin tracer")
	}

	return tracer, nil
}

func neverSample(_ uint64) bool { return false }

type noopRecorder struct{}

func (r *noopRecorder) RecordSpan(zipkintracer.RawSpan) {}

var _ zipkintracer.SpanRecorder = (*noopRecorder)(nil)
