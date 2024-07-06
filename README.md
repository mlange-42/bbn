# BBN

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/bbn/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/bbn/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/bbn)](https://goreportcard.com/report/github.com/mlange-42/bbn)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/bbn)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/bbn)

Bayesian Belief Network CLI/TUI tool and [Go](https://go.dev) module.

![screenshot](https://github.com/mlange-42/bbn/assets/44003176/d81e9225-4480-4e37-a8c0-08ccb02cfe73)

## Features

* Visualize, query and explore networks in the interactive TUI app `bbni`.
* Supports decision networks (aka influence diagrams), including sequential decisions.
* Provides logic nodes for logic inference in addition to probabilistic inference.
* Train and query networks from the command line with `bbn`.
* Human-readable YAML format for networks, as well as BIF-XML.
* Plenty of [examples](https://github.com/mlange-42/bbn/tree/main/_examples) with introductory text, shown in-app.

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

⚠️ Please be aware that the `bbn` Go module is still under development and highly unstable.

Add BBN to a Go project:

```
go get github.com/mlange-42/bbn
```

## Usage

### Command line tools

Try the famous sprinkler example:

```
bbni _examples/bbn/sprinkler.yml
```

Same example with the command line tool, given some evidence:

```
bbn inference _examples/bbn/sprinkler.yml -e Rain=no,GrassWet=yes
```

Train a network from data:

```
bbn train _examples/bbn/fruits-untrained.yml _examples/bbn/fruits.csv
```

Also try the other examples in folder [_examples](https://github.com/mlange-42/bbn/tree/main/_examples).
Run them with `bbni` and play around, but also view their `.yml` files
to get an idea how to create Bayesian Networks.

### Library

⚠️ Please be aware that the `bbn` Go module is still under development and highly unstable.

See the examples in the [API reference](https://pkg.go.dev/github.com/mlange-42/bbn).

## License

This project is distributed under the [MIT license](./LICENSE).
