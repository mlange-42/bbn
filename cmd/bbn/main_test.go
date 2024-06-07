package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	_, _, _, err := run("../../_examples/sprinkler.yml", []string{"Rain=no"}, 100_000, 0)
	assert.Nil(t, err)
}
