package randset

import (
	"sync"
)

// RandomizedSet represents a data structure that stores any comparable type uniquely in a randomized manner.
type RandomizedSet[T comparable] struct {
	m  map[T]struct{}
	mu sync.Mutex
}

// New creates and returns a new instance of RandomizedSet for elements of type T.
func New[T comparable]() *RandomizedSet[T] {
	return &RandomizedSet[T]{
		m: make(map[T]struct{}),
	}
}

// NewWithInitialSize creates and returns a new instance of RandomizedSet with an initial capacity.
func NewWithInitialSize[T comparable](size int) *RandomizedSet[T] {
	return &RandomizedSet[T]{
		m: make(map[T]struct{}, size),
	}
}

// Add adds a new element to the set.
func (s *RandomizedSet[T]) Add(key T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = struct{}{}
}

// Clear removes all elements from the set, resetting its state to empty.
func (s *RandomizedSet[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m = make(map[T]struct{}, len(s.m))
}

// Contains checks whether a given element exists in the set.
func (s *RandomizedSet[T]) Contains(key T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.m[key]

	return exists
}

// Content returns a slice containing all the elements in the set. The order of elements
// in the slice is not specified and can be random.
func (s *RandomizedSet[T]) Content() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := make([]T, 0, len(s.m))

	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

// LoadAndDelete removes and returns a random element from the set. It returns the element
// and true if the set was not empty. If the set is empty, it returns the zero value for T and false.
func (s *RandomizedSet[T]) LoadAndDelete() (T, bool) {
	var (
		key   T
		found bool
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	for key = range s.m {
		found = true

		delete(s.m, key)

		break
	}

	return key, found
}

// IsEmpty checks if the set is empty.
func (s *RandomizedSet[T]) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.m) == 0
}

// Remove removes an element from the set.
func (s *RandomizedSet[T]) Remove(key T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, key)
}

// Size returns the number of elements currently stored in the set.
func (s *RandomizedSet[T]) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.m)
}
