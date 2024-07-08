module github.com/mlange-42/bbn/benchmark

go 1.22.0

require (
	github.com/mlange-42/bbn v0.0.0
	github.com/pkg/profile v1.7.0
)

replace github.com/mlange-42/bbn v0.0.0 => ..

require (
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/google/pprof v0.0.0-20211214055906-6f57359322fd // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
