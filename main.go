package main

import (
	"fmt"

	"github.com/anvaru/gosqldb/btree"
)

func main() {
	bFactor := 3
	b := btree.NewBTree(bFactor)

	keys := []int{
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
		11, 12, 13, 14,
		15, 16, 17, 18,
		19, 20, 21,
	}
	for _, v := range keys {
		b.Insert(v)
		fmt.Println(b.PrettyString())
	}
}
