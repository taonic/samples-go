package query

import (
	"sync"
	"time"

	tdigest "github.com/caio/go-tdigest/v4"
)

type Quantile struct {
	mu sync.Mutex
	t  *tdigest.TDigest
}

func NewQuantile() (*Quantile, error) {
	t, err := tdigest.New()
	if err != nil {
		return nil, err
	}

	return &Quantile{t: t}, nil
}

func (q *Quantile) Add(duration time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.t.Add(float64(duration))
}

func (q *Quantile) Q(p float64) float64 {
	return q.t.Quantile(p)
}
