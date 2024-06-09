# BBN

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/bbn/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/bbn/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/bbn)](https://goreportcard.com/report/github.com/mlange-42/bbn)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/bbn)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/bbn)

Bayesian Belief Network library and CLI/TUI tool for [Go](https://go.dev).

![screenshot](https://github.com/mlange-42/bbn/assets/44003176/0844f5dd-0078-4ba3-8ef8-18441669900a)

## Installation

### Command line tools

As long as there are no precompiled binaries provided, [Go](https://go.dev) is required for installation.

The interactive terminal app `bbni`:

```
go install github.com/mlange-42/bbn/cmd/bbni@latest
```

The command line tool `bbn`:

```
go install github.com/mlange-42/bbn/cmd/bbn@latest
```

### Library

Add BBN to a Go project:

```
go get github.com/mlange-42/bbn
```

## Usage

### Command line tools

Try the famous sprinkler example:

```
bbni _examples/sprinkler.yml
```

Same example with the command line tool, given some evidence:

```
bbn sample _examples/sprinkler.yml -e Rain=no,GrassWet=yes
```

Train a network from data:

```
bbn train _examples/fruits.yml _examples/fruits.csv
```

Also try the other examples in folder [_examples](https://github.com/mlange-42/bbn/tree/main/_examples).

### Library

See the examples in the [API reference](https://pkg.go.dev/github.com/mlange-42/bbn).

## License

This project is distributed under the [MIT license](./LICENSE).
