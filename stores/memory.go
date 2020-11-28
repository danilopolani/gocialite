package stores

import (
	"fmt"
	"sync"
)

type memory struct {
	mu sync.RWMutex
	g  map[string][]byte
}

// NewMemoryStore return a poitner of memory Gocial instance store
func NewMemoryStore() GocialStore {
	return &memory{g: make(map[string][]byte)}
}

func (s *memory) Save(state string, gocial []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.g[state] = gocial
	return nil
}

func (s *memory) Get(state string) ([]byte, error) {
	s.mu.RLock()
	g, ok := s.g[state]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("state not found: %s", state)
	}
	return g, nil
}

func (s *memory) Delete(state string) error {
	s.mu.Lock()
	delete(s.g, state)
	s.mu.Unlock()
	return nil
}

var _ GocialStore = &memory{}
