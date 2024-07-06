# Changelog

## [[unpublished]](https://github.com/mlange-42/bbn/compare/v0.5.0...main)

### Features

* Makes sub-package `ve` for low-level variable elimination public (#73)
* Solve utility nodes even when there are ignored policies, overwritten by "evidence" (#74, #75)

## [[v0.5.0]](https://github.com/mlange-42/bbn/compare/v0.4.0...v0.5.0)

### Features

* Reworked library and apps to use Variable Elimination for inference (#42, #43, #44, #45, #46, #47, #48)
* Implement multi-stage decision networks (#45, #60, #61)
* Adds network training to the TUI app (#54)
* Adds a toggle to the TUI app to ignore policies of decision nodes with evidence (#56)
* Nodes can be moved in the TUI app with WASD keys (#57)
* Save network from TUI app using Ctrl+S (#57)
* Adds a help dialog to the TUI app (#58)
* Optional total utility node for weighting individual utilities (#59)
* Adds support for custom node colors in YAML files (#64)
* Adds shortcuts for logic nodes like, and, or, if-then, etc. (#65)
* Networks can have a longer info text, also shown in TUI app (#66)

### Documentation

* Adds oil drilling decision example (#48)
* Adds robot decision example (#50)
* Adds medical treatment decision example (#60)
* Adds disease control decision example (#61)
* Adds three examples for logic induction using well-known logical puzzles (#65, #66)
* All examples have a detailed description that shows up in the TUI app (#66, #67)

### Other

* Release downloads now include all examples (#68)

## [[v0.4.0]](https://github.com/mlange-42/bbn/compare/v0.3.0...v0.4.0)

### Features

* Uses likelihood-weighted sampling instead of simple rejection sampling (#39)

### Documentation

* Adds an earthquake sensor decision network example (#38, #40)

### Other

* Consistently check whether a new network is an acyclic graph (#40)

## [[v0.3.0]](https://github.com/mlange-42/bbn/compare/v0.2.0...v0.3.0)

### Features

* Adds mouse support for setting outcomes and showing tables (#33)
* Adds utility and decision nodes to form decision networks (#34)
* Color nodes by type (nature, decision, utility) (#35)

## [[v0.2.0]](https://github.com/mlange-42/bbn/compare/v0.1.0...v0.2.0)

### Features

* Support reading BIF-XML file format (#25)

### Documentation

* Adds "Dog Problem" example in BIF-XML format (#25)
* Adds "Native Fish" example in YAML format (#28)
* Adds "Mendel Genetics" example in YAML format (#29)
* Extend Monty-Hall problem example by "Change Door" variable (#30)

### Bugfixes

* Fix outcome labels in probability table view (#27)

## [[v0.1.0]](https://github.com/mlange-42/bbn/commits/v0.1.0/)

Initial release of BBN, the Bayesian Belief Network library and CLI/TUI tool for [Go](https://go.dev).
