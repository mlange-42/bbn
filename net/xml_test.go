package net_test

import (
	"os"
	"testing"

	"github.com/mlange-42/bbn/net"
	"github.com/stretchr/testify/assert"
)

func TestFromBIFXML(t *testing.T) {
	xmlData, err := os.ReadFile("../_examples/dog-problem.xml")
	assert.Nil(t, err)

	net, err := net.FromBIFXML(xmlData)
	assert.Nil(t, err)

	_ = net
}
