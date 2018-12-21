<%=licenseText%>
package state

import "context"

type inMemoryStore struct {
	db map[string]interface{}
}

//NewInMemoryStore New In memory store
func NewInMemoryStore() Store {
	return &inMemoryStore{}
}

func (i *inMemoryStore) Initialize(ctx context.Context) error {
	return nil
}

func (i *inMemoryStore) Close() error {
	return nil
}

func (i *inMemoryStore) Healthy() error {
	return nil
}
