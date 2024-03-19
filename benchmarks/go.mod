module github.com/s3rj1k/go-randset/benchmarks

go 1.22.0

replace github.com/s3rj1k/go-randset => ../.

require (
	github.com/dolthub/swiss v0.2.1
	github.com/s3rj1k/go-randset v0.0.0-00010101000000-000000000000
)

require github.com/dolthub/maphash v0.1.0 // indirect
