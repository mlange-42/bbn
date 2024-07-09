package tui_test

import (
	"testing"

	"github.com/mlange-42/bbn/internal/tui"
	"github.com/stretchr/testify/assert"
)

func TestParseEvidence(t *testing.T) {
	m1, err := tui.ParseEvidence([]string{
		"a=b",
		"c=d",
	})
	assert.Nil(t, err)

	assert.Equal(t, map[string]string{
		"a": "b",
		"c": "d",
	}, m1)

	_, err = tui.ParseEvidence([]string{
		"a=b",
		"c=d=e",
	})
	assert.NotNil(t, err)
}
