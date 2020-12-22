package int64window_test

import (
	"sync"
	"testing"
	"time"

	"github.com/corverroos/ratelimit"
	"github.com/stretchr/testify/require"
)

func TestInt64Window(t *testing.T) {
	l := ratelimit.NewInt64Window(time.Hour, 10)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			require.True(t, l.Request(""))
			wg.Done()
		}()
	}
	wg.Wait()
	require.False(t, l.Request(""))
}

func BenchmarkInt64Window(b *testing.B) {
	ratelimit.Benchmark(b, func() ratelimit.RateLimiter {
		return ratelimit.NewInt64Window(time.Millisecond, 10)
	})
}
