package logic_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/bbn/logic"
	"github.com/stretchr/testify/assert"
)

func TestFactorGet(t *testing.T) {
	f, ok := logic.Get("and")
	assert.True(t, ok)

	table, err := f.Table(2)
	assert.Nil(t, err)
	assert.Equal(t, 8, len(table))

	_, ok = logic.Get("foobar")
	assert.False(t, ok)
}

func Example() {
	and := logic.And()

	table, err := and.Table(2)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(table)/2; i++ {
		fmt.Println(table[i*2 : (i+1)*2])
	}
	// Output:
	//[1 0]
	//[0 1]
	//[0 1]
	//[0 1]
}
