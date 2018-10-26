<%=licenseText%>
package metrics

import (
	"time"
	"context"
	
	"<%=repoUrl%>/pkg/util"
	"<%=repoUrl%>/pkg/log"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// labels are the labels that are send to prometheus
var labels = []string{
	"method",
}

var (
	// StoreCallLatency metric to represent the latency in milliseconds
	storeCallLatency = stats.Float64("store/latency", "The latency in milliseconds per call", "ms")

	// StoreCallCounter metric to represent the number of times store methods are called
	storeCallCount = stats.Int64("store/calls", "number of store calls made", "1")

	// StoreCallErrorCount metric to represent the number of times store methods are called
	StoreCallErrorCount = stats.Int64("store/call_errors", "number of store calls that returned error", "1")
)

var (
	// KeyMethod is the label/tag used while reporting metrics
	KeyMethod, _ = tag.NewKey("method")
)

// NewStoreObserver creates a new StoreObserver
func NewStoreObserver() *StoreObserver {

	storeCallLatencyView := &view.View{
		Name:        "store_call/latency",
		Measure:     storeCallLatency,
		Description: "The distribution of the latencies",

		// Latency in buckets:
		// [>=0ms, >=25ms, >=50ms, >=75ms, >=100ms, >=200ms, >=400ms, >=600ms, >=800ms, >=1s, >=2s, >=4s, >=6s]
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
		TagKeys:     []tag.Key{KeyMethod}}

	storeCallCountView := &view.View{
		Name:        "store_call/count",
		Measure:     storeCallCount,
		Description: "The number calls to the store methods",
		Aggregation: view.Count(),
	}

	storeCallErrorCountView := &view.View{
		Name:        "store_call_error/count",
		Measure:     StoreCallErrorCount,
		Description: "The number store calls which returnd in error to the store methods",
		Aggregation: view.Count(),
	}

	log.Debug("Registering store metrics views...")
	// Register the views
	if err := view.Register(storeCallLatencyView, storeCallCountView, storeCallErrorCountView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}

	return &StoreObserver{}
}

// StoreObserver encapsulates exposing of store specific metrics to Prometheus.
type StoreObserver struct {
}

// defaultIgnoredMethods are methods which are commonly found on our stores and
// thus ignored when preloading.
var defaultIgnoredMethods = []string{"Close", "Healthy", "C"}

// Preload counters and histograms for each method defined on s. You can
// optionally supply extra ignoreMethods which will be added to the
// defaultIgnoredMethods array.
func (s *StoreObserver) Preload(ifc interface{}, extraIgnoredMethods ...string) {
	ignoredMethods := append(defaultIgnoredMethods, extraIgnoredMethods...)
	methods := util.GetMethods(ifc)
	for _, method := range methods {
		if shouldIgnore(method, ignoredMethods) {
			continue
		}

	}
}

func shouldIgnore(method string, ignoredMethods []string) bool {
	for _, ignore := range ignoredMethods {
		if method == ignore {
			return true
		}
	}

	return false
}

// Observe immediately increments the counter for method and returns a func
// which will observe an metric item in duration based on the duration
func (s *StoreObserver) Observe(ctx context.Context, method string) func() {
	ctx, err := tag.New(ctx, tag.Insert(KeyMethod, method))
	if err != nil {
		log.Fatalf("Failed to Observe method %s: %v", method, err)
	}

	stats.Record(ctx, storeCallCount.M(1)) // Counter to track a store call
	startTime := time.Now()

	return func() {
		ms := float64(time.Since(startTime).Nanoseconds()) / 1e6
		stats.Record(ctx, storeCallLatency.M(ms))
	}
}
