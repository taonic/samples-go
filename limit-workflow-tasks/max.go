package limit

import (
	"sync"
)

type Max struct {
	mu sync.Mutex
	M  int
}

func (m *Max) Compare(value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.M < value {
		m.M = value
	}
}
