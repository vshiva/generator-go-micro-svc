# metrics

Metrics contains metrics exposing and collecting related code.

## StoreObserver

The [`StoreObserver`](store_observer.go) encapsulates exposing Store specific
metrics to Prometheus. The following metrics are exposed:

- `store_handling_seconds` a histogram containing the response latency of the
  calls.
- `store_handled_total` a counter which keeps track of all call made.

A specific metrics store should instantiate this StoreObserver and proxy all
calls to `Observe`. This method returns a func() which has to be called once
the metrics store has finished the call:

```
func (s *MetricsStore) GetUser(ctx context.Context, id string) (User, error) {
  done := s.observer.Observe("GetUser")
  defer done()
  // ... actual work
  return user, err
}
```

It is highly recommended to call `Preload` with the store to expose all metrics
to prometheus from the beginning:

```
func NewMetricsStore(wrappedStore Store) *MetricsStore {
	store := &MetricsStore{store: wrappedStore, observer: metrics.NewStoreObserver()}
	store.observer.Preload(store)
	return store
}
```
