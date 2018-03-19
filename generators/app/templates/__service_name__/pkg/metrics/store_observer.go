<%=licenseText%>
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"<%=repoUrl%>/pkg/util"
)

// labels are the labels that are send to prometheus
var labels = []string{
	"method",
}

// NewStoreObserver creates a new StoreObserver
func NewStoreObserver() *StoreObserver {
	durationOpts := prometheus.HistogramOpts{
		Name: "store_handling_seconds",
		Help: "Histogram of response latency (seconds) of store calls that had been handled by the server",
	}
	duration := prometheus.NewHistogramVec(durationOpts, labels)

	counterOpts := prometheus.CounterOpts{
		Name: "store_handled_total",
		Help: "Total number of store calls completed on the server, regardless of success or failure",
	}
	counter := prometheus.NewCounterVec(counterOpts, labels)

	prometheus.MustRegister(duration)
	prometheus.MustRegister(counter)

	return &StoreObserver{duration: duration, counter: counter}
}

// StoreObserver encapsulates exposing of store specific metrics to Prometheus.
type StoreObserver struct {
	duration *prometheus.HistogramVec
	counter  *prometheus.CounterVec
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

		s.counter.WithLabelValues(method)
		s.duration.WithLabelValues(method)
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
// which will observe an metric item in duration based on the duration.
func (s *StoreObserver) Observe(method string) func() {
	start := time.Now()

	counter := s.counter.WithLabelValues(method)
	counter.Add(1)

	duration := s.duration.WithLabelValues(method)
	return func() {
		duration.Observe(time.Now().Sub(start).Seconds())
	}
}
