package logic_test

import (
	"fmt"

	"github.com/mlange-42/bbn/logic"
)

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
