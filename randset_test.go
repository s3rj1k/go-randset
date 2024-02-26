package randset_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sync"
	"testing"

	"github.com/s3rj1k/go-randset"
)

const (
	fnAdd           = "Add"
	fnClear         = "Clear"
	fnContains      = "Contains"
	fnContent       = "Content"
	fnIsEmpty       = "IsEmpty"
	fnLoadAndDelete = "LoadAndDelete"
	fnRemove        = "Remove"
	fnSize          = "Size"
)

func TestRandomizedSet(t *testing.T) {
	tests := []struct {
		name            string
		operations      []string
		values          [][]int
		expectedResults []any
	}{
		{
			name:            fnIsEmpty,
			operations:      []string{fnIsEmpty},
			values:          [][]int{{}},
			expectedResults: []any{true},
		},
		{
			name:            fmt.Sprintf("%s and %s", fnAdd, fnContains),
			operations:      []string{fnAdd, fnContains},
			values:          [][]int{{1}, {1}},
			expectedResults: []any{nil, true},
		},
		{
			name:            fmt.Sprintf("%s %s and %s", fnAdd, fnRemove, fnContains),
			operations:      []string{fnAdd, fnRemove, fnContains},
			values:          [][]int{{2}, {2}, {2}},
			expectedResults: []any{nil, nil, false},
		},
		{
			name:            fmt.Sprintf("%s and %s", fnClear, fnIsEmpty),
			operations:      []string{fnAdd, fnContains, fnClear, fnIsEmpty},
			values:          [][]int{{1}, {1}, {}, {}},
			expectedResults: []any{nil, true, nil, true},
		},
		{
			name:            fmt.Sprintf("%s and %s", fnContent, fnSize),
			operations:      []string{fnAdd, fnAdd, fnContent, fnSize},
			values:          [][]int{{1}, {2}, {}, {}},
			expectedResults: []any{nil, nil, []int{1, 2}, 2},
		},
		{
			name:            fnLoadAndDelete,
			operations:      []string{fnAdd, fnLoadAndDelete, fnContains, fnIsEmpty},
			values:          [][]int{{1}, {}, {1}, {}},
			expectedResults: []any{nil, true, false, true},
		},
	}

	for i := range tests {
		tt := tests[i]

		t.Run(tt.name, func(t *testing.T) {
			set := randset.New[int]()

			for i, operation := range tt.operations {
				switch operation {
				case fnAdd:
					set.Add(tt.values[i][0])
				case fnClear:
					set.Clear()
				case fnContains:
					got := set.Contains(tt.values[i][0])
					if got != tt.expectedResults[i] {
						t.Errorf("%s: %s(%d) = %v, want %v",
							tt.name, operation, tt.values[i][0], got, tt.expectedResults[i],
						)
					}
				case fnContent:
					got := set.Content()
					slices.Sort(got)

					expected, ok := tt.expectedResults[i].([]int)
					if !ok {
						t.Fatalf("Unexpected value type: %T", expected)
					}

					slices.Sort(expected)

					if !reflect.DeepEqual(got, expected) {
						t.Errorf("%s: %s() = %v, want %v", tt.name, operation, got, expected)
					}
				case fnLoadAndDelete:
					removedKey, removed := set.LoadAndDelete()
					if !removed || set.Contains(removedKey) {
						t.Errorf("%s: %s() failed", tt.name, operation)
					}
				case fnIsEmpty:
					got := set.IsEmpty()
					if got != tt.expectedResults[i] {
						t.Errorf("%s: %s() = %v, want %v", tt.name, operation, got, tt.expectedResults[i])
					}
				case fnRemove:
					set.Remove(tt.values[i][0])
				case fnSize:
					got := set.Size()
					if got != tt.expectedResults[i] {
						t.Errorf("%s: %s() = %v, want %v", tt.name, operation, got, tt.expectedResults[i])
					}
				}
			}
		})
	}
}

const (
	setSize = 500000
)

func TestLoadAndDelete(t *testing.T) {
	set := randset.NewWithInitialSize[uint64](setSize)

	for set.Size() < setSize {
		set.Add(rand.Uint64()) //nolint:gosec // G404: Use of weak random number generator
	}

	if set.Size() != setSize {
		t.Fatalf("Initial setup failed, set size = %d, expected %d", set.Size(), setSize)
	}

	for i := range setSize {
		beforeSize := set.Size()

		val, found := set.LoadAndDelete()
		if !found {
			t.Fatalf("%s() failed to return a value on iteration %d", fnLoadAndDelete, i)
		}

		if set.Contains(val) {
			t.Errorf("Expected the value %d to be removed from set, but it was not", val)
		}

		afterSize := set.Size()

		if beforeSize-1 != afterSize {
			t.Errorf("Expected size to decrease by 1, but it went from %d to %d on iteration %d", beforeSize, afterSize, i)
		}
	}

	if !set.IsEmpty() {
		t.Errorf("Expected the set to be empty after %d removals, but it was not", setSize)
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
			b.Fatalf("%s() failed to return a value, set might be empty", fnLoadAndDelete)
		}

		set.Add(val)
	}
}

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
