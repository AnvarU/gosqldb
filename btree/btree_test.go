package btree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBtreeCreation(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)
	assert.Equal(t, b.GetMaxKeys(), 2*bFactor-1)
	assert.Equal(t, b.GetMaxChild(), 2*bFactor)
}

func TestBtreeInsertWihoutSplit(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	b.Insert(50)
	assert.Equal(t, len(b.Root.Keys), 1)
	assert.Equal(t, b.Root.Keys[0], 50)

	b.Insert(30)
	b.Insert(70)
	b.Insert(40)

	assert.Equal(t, len(b.Root.Keys), 4)
	assert.Equal(t, b.Root.Keys, []int{30, 40, 50, 70})
}

func TestBtreeSplittingRoot(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	b.Insert(50)
	assert.Equal(t, len(b.Root.Keys), 1)
	assert.Equal(t, b.Root.Keys[0], 50)

	b.Insert(30)
	b.Insert(70)
	b.Insert(40)
	assert.Equal(t, len(b.Root.Keys), 4)
	assert.Equal(t, b.Root.Keys, []int{30, 40, 50, 70})

	b.Insert(35)
	b.Insert(32)
	assert.Equal(t, len(b.Root.Keys), 1)
	assert.Equal(t, len(b.Root.Children), 2)
	assert.Equal(t, b.Root.Keys[0], 40)
	assert.Equal(t, b.Root.Children[0].Keys, []int{30, 32, 35})
	assert.Equal(t, b.Root.Children[1].Keys, []int{50, 70})
}

func TestBtreeInsertWithSplit(t *testing.T) {
	bFactor := 3
	b := NewBTree(bFactor)

	b.Insert(50)
	b.Insert(30)
	b.Insert(70)
	b.Insert(40)
	b.Insert(35)
	b.Insert(32)
	b.Insert(31)
	b.Insert(33)
	b.Insert(34)
	b.Insert(80)
	b.Insert(90)
	b.Insert(100)
	b.Insert(110)
	b.Insert(120)
	b.Insert(130)
	b.Insert(140)
	b.Insert(150)
	b.Insert(160)
	b.Insert(170)
	b.Insert(180)
	assert.Equal(t, len(b.Root.Keys), 1)
	assert.Equal(t, len(b.Root.Children), 2)
	assert.Equal(t, b.Root.Keys, []int{80})
	assert.Equal(t, b.Root.Children[0].Keys, []int{32, 40})
	assert.Equal(t, len(b.Root.Children[0].Children), 3)
	assert.Equal(t, b.Root.Children[0].Children[0].Keys, []int{30, 31})
	assert.Equal(t, b.Root.Children[0].Children[1].Keys, []int{33, 34, 35})
	assert.Equal(t, b.Root.Children[0].Children[2].Keys, []int{50, 70})

	assert.Equal(t, b.Root.Children[1].Keys, []int{110, 140})
	assert.Equal(t, len(b.Root.Children[1].Children), 3)
	assert.Equal(t, b.Root.Children[1].Children[0].Keys, []int{90, 100})
	assert.Equal(t, b.Root.Children[1].Children[1].Keys, []int{120, 130})
	assert.Equal(t, b.Root.Children[1].Children[2].Keys, []int{150, 160, 170, 180})
}
