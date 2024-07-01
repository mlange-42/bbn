package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunInferenceCommand(t *testing.T) {
	_, _, _, err := runInferenceCommand("../../_examples/sprinkler.yml", []string{"Rain=no"})
	assert.Nil(t, err)
}
