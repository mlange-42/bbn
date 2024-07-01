# BBN

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/bbn/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/bbn/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/bbn)](https://goreportcard.com/report/github.com/mlange-42/bbn)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/bbn)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/bbn)

Bayesian Belief Network library and CLI/TUI tool for [Go](https://go.dev).

![screenshot](https://github.com/mlange-42/bbn/assets/44003176/0844f5dd-0078-4ba3-8ef8-18441669900a)

## Features

* Minimal, fast API for usage as a library.
* Decision networks with utility and decision nodes.
* Human-readable YAML format, as well as BIF-XML.
* Train and query networks from the command line with `bbn`.
* Visualize, query and explore networks in the interactive TUI app `bbni`.

## Installation

### Command line tools

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/bbn/releases).

> Alternatively, install the latest development versions of `bbn` and `bbni` using [Go](https://go.dev):
> ```shell
> go install github.com/mlange-42/bbn/cmd/bbn@main
> go install github.com/mlange-42/bbn/cmd/bbni@main
> ```

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
bbn inference _examples/sprinkler.yml -e Rain=no,GrassWet=yes
```

~~Train a network from data:~~ (currently not functional)

```
bbn train _examples/fruits-untrained.yml _examples/fruits.csv
```

Also try the other examples in folder [_examples](https://github.com/mlange-42/bbn/tree/main/_examples).

### Library

See the examples in the [API reference](https://pkg.go.dev/github.com/mlange-42/bbn).

## License

This project is distributed under the [MIT license](./LICENSE).
