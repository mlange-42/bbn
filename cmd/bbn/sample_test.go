package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSampleCommand(t *testing.T) {
	_, _, _, err := runSampleCommand("../../_examples/sprinkler.yml", []string{"Rain=no"}, 100_000, 0)
	assert.Nil(t, err)
}
