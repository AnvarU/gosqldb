package main

import "github.com/anvaru/gosqldb/btree"

func main() {
	b := btree.NewBTree(3)
	b.Insert(30)
	b.Insert(40)
	b.Insert(20)
	b.Insert(60)
	b.Insert(50)
	b.Insert(10)
	b.Insert(70)
	b.Insert(33)
}
