# BBN

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/bbn/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/bbn/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/bbn)](https://goreportcard.com/report/github.com/mlange-42/bbn)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/bbn)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/bbn)

Bayesian Belief Network library and CLI tool for [Go](https://go.dev).

## Installation

### Command line tool

As long as there are no precompiled binaries provided, [Go](https://go.dev) is required for installation:

```
go install github.com/mlange-42/bbn/cmd/bbn@latest
```

### Library

Add BBN to a Go project:

```
go get github.com/mlange-42/bbn
```

## Usage

### Command line tool

Try the sprinkler example:

```
bbn _examples/sprinkler.yml
```

The same example, given some evidence:

```
bbn _examples/sprinkler.yml -e Rain=no,GrassWet=yes
```

### Library

See the examples in the [API reference](https://pkg.go.dev/github.com/mlange-42/bbn).

## License

This project is distributed under the [MIT license](./LICENSE).
