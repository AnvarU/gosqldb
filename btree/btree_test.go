package btree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBtreeCreation(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)
	assert.Equal(t, 2*bFactor-1, b.GetMaxKeys())
	assert.Equal(t, 2*bFactor, b.GetMaxChild())
}

func TestBtreeInsertWihoutSplit(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	b.Insert(50)
	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, 50, b.Root.Keys[0])

	b.Insert(30)
	b.Insert(70)
	b.Insert(40)

	assert.Equal(t, 4, len(b.Root.Keys))
	assert.Equal(t, []int{30, 40, 50, 70}, b.Root.Keys)
}

func TestBtreeSplittingRoot(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	b.Insert(50)
	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, 50, b.Root.Keys[0])

	b.Insert(30)
	b.Insert(70)
	b.Insert(40)
	assert.Equal(t, 4, len(b.Root.Keys))
	assert.Equal(t, []int{30, 40, 50, 70}, b.Root.Keys)

	b.Insert(35)
	b.Insert(32)
	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, 2, len(b.Root.Children))
	assert.Equal(t, 40, b.Root.Keys[0])
	assert.Equal(t, []int{30, 32, 35}, b.Root.Children[0].Keys)
	assert.Equal(t, []int{50, 70}, b.Root.Children[1].Keys)
}

func TestBtreeInsertWithSplit(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	keys := []int{
		50, 30, 70, 40, 35,
		32, 31, 33, 34, 80,
		90, 100, 110, 120,
		130, 140, 150, 160,
		170, 180,
	}
	for _, v := range keys {
		b.Insert(v)
	}

	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, 2, len(b.Root.Children))
	assert.Equal(t, []int{80}, b.Root.Keys)
	assert.Equal(t, []int{32, 40}, b.Root.Children[0].Keys)
	assert.Equal(t, 3, len(b.Root.Children[0].Children))
	assert.Equal(t, []int{30, 31}, b.Root.Children[0].Children[0].Keys)
	assert.Equal(t, []int{33, 34, 35}, b.Root.Children[0].Children[1].Keys)
	assert.Equal(t, []int{50, 70}, b.Root.Children[0].Children[2].Keys)

	assert.Equal(t, []int{110, 140}, b.Root.Children[1].Keys)
	assert.Equal(t, 3, len(b.Root.Children[1].Children))
	assert.Equal(t, []int{90, 100}, b.Root.Children[1].Children[0].Keys)
	assert.Equal(t, []int{120, 130}, b.Root.Children[1].Children[1].Keys)
	assert.Equal(t, []int{150, 160, 170, 180}, b.Root.Children[1].Children[2].Keys)
}

func TestBtreeDeletionFromLeaf(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	keys := []int{
		1, 2, 3, 4, 5, 6,
	}
	for _, v := range keys {
		b.Insert(v)
	}

	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, []int{4, 5, 6}, b.Root.Children[1].Keys)

	b.Delete(5)
	assert.Equal(t, []int{4, 6}, b.Root.Children[1].Keys)
}

func TestBtreeDeletionFromInternalNode(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	keys := []int{
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
		11, 12, 13, 14,
		15, 16, 17, 18,
		19, 20, 21,
	}
	for _, v := range keys {
		b.Insert(v)
	}

	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, []int{12, 15, 18}, b.Root.Children[1].Keys)

	b.Delete(18)
	assert.Equal(t, []int{12, 15, 19}, b.Root.Children[1].Keys)
}

func TestBtreeDeletionWithMerging(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	keys := []int{
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
		11, 12, 13, 14,
		15, 16, 17, 18,
		19, 20, 21,
	}
	for _, v := range keys {
		b.Insert(v)
	}

	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, []int{12, 15, 18}, b.Root.Children[1].Keys)

	b.Delete(15)
	assert.Equal(t, []int{12, 18}, b.Root.Children[1].Keys)
	assert.Equal(t, []int{13, 14, 16, 17}, b.Root.Children[1].Children[1].Keys)
}

func TestBtreeDeletionWithFixingNodes(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	keys := []int{
		50, 30, 70, 40, 35,
		32, 31, 33, 34, 80,
		90, 100, 110, 120,
		130, 140, 150, 160,
		170, 180,
	}
	for _, v := range keys {
		b.Insert(v)
	}

	assert.Equal(t, 1, len(b.Root.Keys))
	assert.Equal(t, 2, len(b.Root.Children))
	assert.Equal(t, []int{80}, b.Root.Keys)
	assert.Equal(t, []int{32, 40}, b.Root.Children[0].Keys)
	assert.Equal(t, 3, len(b.Root.Children[0].Children))
	assert.Equal(t, []int{30, 31}, b.Root.Children[0].Children[0].Keys)
	assert.Equal(t, []int{33, 34, 35}, b.Root.Children[0].Children[1].Keys)
	assert.Equal(t, []int{50, 70}, b.Root.Children[0].Children[2].Keys)

	assert.Equal(t, []int{110, 140}, b.Root.Children[1].Keys)
	assert.Equal(t, 3, len(b.Root.Children[1].Children))
	assert.Equal(t, []int{90, 100}, b.Root.Children[1].Children[0].Keys)
	assert.Equal(t, []int{120, 130}, b.Root.Children[1].Children[1].Keys)
	assert.Equal(t, []int{150, 160, 170, 180}, b.Root.Children[1].Children[2].Keys)

	b.Delete(34)
	assert.Equal(t, []int{33, 35}, b.Root.Children[0].Children[1].Keys)
}
