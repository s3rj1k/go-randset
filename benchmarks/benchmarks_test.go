package benchmarks_test

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/dolthub/swiss"
	"github.com/s3rj1k/go-randset"
)

const (
	setSize = 500000
)

func BenchmarkSliceQueue(b *testing.B) {
	s := make([]uint64, setSize)

	for i := 0; i < setSize; i++ {
		s[i] = rand.Uint64() //nolint:gosec // G404: Use of weak random number generator
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val := s[0]

		s = s[1:]
		s = append(s, val)
	}
}

func BenchmarkChannelQueue(b *testing.B) {
	ch := make(chan uint64, setSize)

	for range setSize {
		ch <- rand.Uint64() //nolint:gosec // G404: Use of weak random number generator
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val := <-ch
		ch <- val
	}
}

func BenchmarkSwissTableQueue(b *testing.B) {
	sm := swiss.NewMap[uint64, struct{}](setSize)

	for range setSize {
		for {
			key := rand.Uint64() //nolint:gosec // G404: Use of weak random number generator

			if found := sm.Has(key); found {
				continue
			}

			sm.Put(key, struct{}{})

			break
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var key uint64

		sm.Iter(func(k uint64, _ struct{}) (stop bool) {
			key = k

			if ok := sm.Delete(k); !ok {
				b.Fatal("Failed to delete data from SwissTable")
			}

			return true
		})

		sm.Put(key, struct{}{})
	}
}

func BenchmarkRandomSet(b *testing.B) {
	set := randset.NewWithInitialSize[uint64](setSize)

	for set.Size() < setSize {
		set.Add(rand.Uint64()) //nolint:gosec // G404: Use of weak random number generator
	}

	if set.Size() != setSize {
		b.Fatalf("Initial setup failed, set size = %d, expected %d", set.Size(), setSize)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val, found := set.LoadAndDelete()
		if !found {
			b.Fatal("Failed to return a value, set might be empty")
		}

		set.Add(val)
	}
}

func BenchmarkSyncMapQueue(b *testing.B) {
	var sm sync.Map

	for range setSize {
		for {
			key := rand.Uint64() //nolint:gosec // G404: Use of weak random number generator

			_, found := sm.Load(key)
			if found {
				continue
			}

			sm.Store(key, struct{}{})

			break
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var (
			key uint64
			ok  bool
		)

		sm.Range(func(k, _ any) bool {
			key, ok = k.(uint64)
			if !ok {
				b.Fatalf("Unexpected data type: %T", key)
			}

			return false
		})

		_, ok = sm.LoadAndDelete(key)
		if !ok {
			b.Fatal("Failed to load data from sync.Map")
		}

		sm.Store(key, struct{}{})
	}
}
