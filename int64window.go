package int64window

import (
	"sync"
	"time"

	"github.com/corverroos/ratelimit"
)

func NewInt64Window(period time.Duration, limit int) *RateLimiter {
	w := &Int64Window{
		period: period,
		limit:  limit,
		counts: make(map[string]counter, 0),
	}

	go w.CurrentStop()

	return w
}

type counter struct {
	count   int
	limited bool
}

func (c *counter) UpdateCounter(l int) bool {
	if c.limited {
		return false
	}

	c.count++
	if c.count <= l {
		return true
	}
	c.limited = true

	return false
}

func (n *Int64Window) CurrentStop() {
	for {
		time.Sleep(n.period)
		n.mu.Lock()
		n.counts = make(map[string]counter, 0)
		n.mu.Unlock()
	}
}

type Int64Window struct {
	bucketSize int64
	limit      int
	period     time.Duration

	currentStop int64
	counts      map[string]counter
	mu          sync.Mutex
}

func (n *Int64Window) Request(resource string) bool {
	n.mu.Lock()
	i := n.counts[resource]
	cmp := i.UpdateCounter(n.limit)
	n.mu.Unlock()
	return cmp
}

var _ ratelimit.RateLimiter = (*Int64Window)(nil)
