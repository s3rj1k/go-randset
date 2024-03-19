# RandSet

`RandSet` is a Go module offering an arguably efficient and concurrency-safe implementation of a set data structure, uniquely designed to support random element eviction. This module is particularly suited for applications requiring a queue-like data structure without the need for maintaining queue order, yet expecting the element uniqueness characteristic of a set.

## Features

- **Efficient Operations**: Constant time complexity for basic methods.
- **Random Element Access**: Retrieve and Evict random element from the set in constant time.
- **Concurrency Safe**: Support for concurrent access with thread-safe wrappers.
- **Zero Dependencies**: Uses only Go's standard library, ensuring easy integration and compatibility.

## Benchmark results.
`$ go test -run=^$ -bench=. -benchmem`

```
goos: linux
goarch: amd64
pkg: github.com/s3rj1k/go-randset
cpu: 12th Gen Intel(R) Core(TM) i5-1240P
BenchmarkSliceQueue-16          226308338          4.814 ns/op    39 B/op        0 allocs/op
BenchmarkChannelQueue-16         40899385         29.27 ns/op      0 B/op        0 allocs/op
BenchmarkSwissTableQueue-16      20348682        381.9 ns/op       3 B/op        0 allocs/op
BenchmarkRandomSet-16            12879373         88.29 ns/op      0 B/op        0 allocs/op
BenchmarkSyncMapQueue-16          3520684        335.1 ns/op      24 B/op        2 allocs/op
```

## Getting Started

### Prerequisites

- Go 1.22.0 or later.

### Installation

To install `RandSet`, use the `go get` command:

```bash
go get github.com/s3rj1k/go-randset
```

### Example

```
package main

import (
	"fmt"

	"github.com/s3rj1k/go-randset"
)

func main() {
	// Create a new RandomizedSet of int type.
	set := randset.New[int]()

	// Add elements.
	for i := range 5 {
		set.Add(i)
	}

	// Get the content of the set.
	fmt.Printf("Content of the set: %v\n", set.Content())

	// Check if the set contains a specific element.
	fmt.Printf("Set contains 2: %t\n", set.Contains(2))

	// Check size.
	fmt.Printf("Size of the set: %d\n", set.Size())

	// Remove an element.
	set.Remove(2)
	fmt.Printf("After removing 2, size of the set: %d\n", set.Size())

	// Check if set is empty.
	fmt.Printf("Set empty: %t\n", set.IsEmpty())

	// Load and Delete an element.
	val, ok := set.LoadAndDelete()
	if ok {
		fmt.Printf("Loaded and deleted: %d\n", val)
		fmt.Printf("Set contains %d: %t\n", val, set.Contains(val))
	}

	// Clear the set.
	set.Clear()
	fmt.Printf("After clear, set is empty: %t\n", set.IsEmpty())
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue to discuss potential improvements or features.

## License

`RandSet` is available under the MIT license. See the LICENSE file for more info.
