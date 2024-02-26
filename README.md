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
BenchmarkRandomSet-16            15442075        83.31 ns/op         0 B/op        0 allocs/op
BenchmarkSliceQueue-16          231713943        5.094 ns/op        39 B/op        0 allocs/op
BenchmarkChannelQueue-16         40679883        29.27 ns/op         0 B/op        0 allocs/op
BenchmarkSyncMapQueue-16          3498108        338.0 ns/op        24 B/op        2 allocs/op
PASS
ok  	github.com/s3rj1k/go-randset	6.887s
```

## Getting Started

### Prerequisites

- Go 1.22.0 or later.

### Installation

To install `RandSet`, use the `go get` command:

```bash
go get github.com/s3rj1k/go-randset
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue to discuss potential improvements or features.

## License

`RandSet` is available under the MIT license. See the LICENSE file for more info.
