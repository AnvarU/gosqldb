package main

import (
	"fmt"

	"github.com/anvaru/gosqldb/btree"
)

func main() {
	bFactor := 3
	b := btree.NewBTree(bFactor)

	for i := 1; i < 69; i++ {
		b.Insert(i)
	}

	fmt.Println(b.PrettyString())
}
