package bbn_test

import (
	"os"
	"testing"

	"github.com/mlange-42/bbn"
	"github.com/stretchr/testify/assert"
)

func TestFromBIFXML(t *testing.T) {
	xmlData, err := os.ReadFile("_examples/dog-problem.xml")
	assert.Nil(t, err)

	net, err := bbn.FromBIFXML(xmlData)
	assert.Nil(t, err)

	_ = net
}
