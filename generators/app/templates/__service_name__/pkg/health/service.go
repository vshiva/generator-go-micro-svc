<%=licenseText%>
package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"<%=repoUrl%>/pkg/log"
)

// New create a new HealthService
func New() *Service {
	return &Service{
		probes: make(map[string]Probe),
	}
}

// The Service periodiablly checks all probes, if any fails it will keep
// track of this. It exposes this information through two handlers,
// livenesProbe and readinessProbe.
type Service struct {
	mu     sync.Mutex
	probes map[string]Probe

	server *http.Server

	healthyCount   int
	unhealthyCount int
}

// ListenAndServe starts looping through the probes and it will start the
// http server on addr.
func (s *Service) ListenAndServe(addr string) error {
	go s.healthcheck()

	m := http.NewServeMux()

	m.HandleFunc("/live", s.livenessProbe)
	m.HandleFunc("/ready", s.readinessProbe)

	s.server = &http.Server{
		Addr:    addr,
		Handler: m,
	}
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open
// listeners, then closing all idle connections, and then waiting
// indefinitely for connections to return to idle and then shut down.
// If the provided context expires before the shutdown is complete,
// then the context's error is returned.
// This operation is a no-op if no server has been started.
func (s *Service) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// healthcheck starts an infinite loop which will iterate over all probes to
// see if there are unhealthy probes.
func (s *Service) healthcheck() {
	// TODO: Support stopping this infinite loop
	for {
		healthy := true

		// iterate all probes
		s.mu.Lock()
		for name, probe := range s.probes {
			// TODO: Add
			err := probe.Healthy()
			if err != nil {
				healthy = false
				log.Warn("Probe %s failed its healthcheck: %+v", name, err)
			}
		}

		sleepDuration := time.Second * 5

		if healthy {
			s.healthyCount++
			s.unhealthyCount = 0
		} else {
			s.unhealthyCount++
			s.healthyCount = 0
			sleepDuration = time.Second * 2
		}

		s.mu.Unlock()

		time.Sleep(sleepDuration)
	}
}

// livenessProbe reports bad health when a probe failed 5 times. It should be
// used to terminate this service.
func (s *Service) livenessProbe(res http.ResponseWriter, req *http.Request) {
	if s.unhealthyCount > 5 {
		http.Error(res, "500 bad health", http.StatusInternalServerError)
		return
	}
}

// readynessProbe reports bad health when a probe failed 1 time. It should be
// used to temporary prevent traffic from coming to this service.
func (s *Service) readinessProbe(res http.ResponseWriter, req *http.Request) {
	if s.unhealthyCount > 0 {
		http.Error(res, "500 bad health", http.StatusInternalServerError)
		return
	}
}

// RegisterProbe adds a new probe to be monitored.
func (s *Service) RegisterProbe(name string, p Probe) {
	s.mu.Lock()

	s.probes[name] = p

	s.mu.Unlock()
}
