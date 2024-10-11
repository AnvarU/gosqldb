package btree

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type Node struct {
	Keys       []int
	Children   []*Node
	parentNode *Node
}

func (n *Node) String() string {
	return fmt.Sprintf("%v ", n.Keys)
}

func (n *Node) splitChild(i int) {
	childNode := n.Children[i]
	leftNode := &Node{
		Keys: make([]int, len(childNode.Keys)/2),
	}
	copy(leftNode.Keys, childNode.Keys[:len(childNode.Keys)/2])

	rightNode := &Node{
		Keys: make([]int, len(childNode.Keys)/2),
	}
	copy(rightNode.Keys, childNode.Keys[len(childNode.Keys)/2+1:])

	if len(childNode.Children) != 0 {

		for k, v := range childNode.Children {
			if k < len(childNode.Children)/2 {
				leftNode.Children = append(leftNode.Children, v)
			} else {
				rightNode.Children = append(rightNode.Children, v)
			}
		}

	}

	leftoverKey := childNode.Keys[len(childNode.Keys)/2]
	n.Keys = append(n.Keys, 0)
	copy(n.Keys[i+1:], n.Keys[i:])
	n.Keys[i] = leftoverKey

	n.Children[i] = leftNode
	i += 1
	n.Children = append(n.Children, &Node{})
	copy(n.Children[i+1:], n.Children[i:])
	n.Children[i] = rightNode
}

func (n *Node) insert(b *BTree, k int) *Node {
	if n == nil {
		return &Node{
			Keys:     []int{k},
			Children: make([]*Node, 0),
		}
	}

	index := sort.Search(len(n.Keys), func(i int) bool {
		return n.Keys[i] > k
	})
	if len(n.Children) == 0 {
		n.Keys = append(n.Keys, 0)
		copy(n.Keys[index+1:], n.Keys[index:])
		n.Keys[index] = k
	} else {
		childNode := n.Children[index]
		if len(childNode.Keys) == b.GetMaxKeys() {
			n.splitChild(index)

			if k > n.Keys[index] {
				index += 1
			}
			childNode = n.Children[index]
		}
		childNode.insert(b, k)
	}

	return n
}

func (n *Node) deleteFromLeftChild(b *BTree, i int) int {
	child := n.Children[i]
	keyCount := len(child.Keys)
	if len(child.Children) == 0 {
		key := child.Keys[keyCount-1]
		child.Keys = child.Keys[:keyCount-1]
		return key
	}
	cCount := len(child.Children)
	if len(child.Children[cCount-1].Keys) > b.GetMinKeys() {
		return child.deleteFromLeftChild(b, cCount-1)
	} else if len(child.Children[cCount-2].Keys) > b.GetMinKeys() {
		child.Children[cCount-1].Keys = append(child.Children[cCount-1].Keys, 0)
		copy(child.Children[cCount-1].Keys[1:], child.Children[cCount-1].Keys[:len(child.Children[cCount-1].Keys)-1])
		child.Children[cCount-1].Keys[0] = child.Keys[keyCount-1]

		child.Keys[keyCount-1] = child.Children[cCount-2].Keys[len(child.Children[cCount-2].Keys)-1]
		child.Children[cCount-2].Keys = child.Children[cCount-2].Keys[:len(child.Children[cCount-2].Keys)-1]
		return child.deleteFromLeftChild(b, cCount-1)
	} else {
		child.mergeChildren(cCount - 2)
		return child.deleteFromLeftChild(b, cCount-2)
	}
}

func (n *Node) deleteFromRightChild(b *BTree, i int) int {
	child := n.Children[i]
	if len(child.Children) == 0 {
		key := child.Keys[0]
		child.Keys = child.Keys[1:]
		return key
	}
	if len(child.Children[0].Keys) > b.GetMinKeys() {
		return child.deleteFromRightChild(b, 0)
	} else if len(child.Children[1].Keys) > b.GetMinKeys() {
		child.Children[0].Keys = append(child.Children[0].Keys, child.Keys[0])
		child.Keys[0] = child.Children[1].Keys[0]
		child.Children[1].Keys = child.Children[1].Keys[1:]
		return child.deleteFromRightChild(b, 0)
	} else {
		child.mergeChildren(0)
		return child.deleteFromRightChild(b, 0)
	}
}

func (n *Node) mergeChildren(i int) {
	lchild := n.Children[i]
	rchild := n.Children[i+1]

	lchild.Keys = append(lchild.Keys, n.Keys[i])
	lchild.Keys = append(lchild.Keys, rchild.Keys...)
	lchild.Children = append(lchild.Children, rchild.Children...)

	n.Keys = append(n.Keys[:i], n.Keys[i+1:]...)
	n.Children = append(n.Children[:i+1], n.Children[i+2:]...)
}

func (n *Node) deleteAndMergeChildren(b *BTree, k, i int) {
	n.mergeChildren(i)
	n.Children[i].delete(b, k)
}

func (n *Node) deleteInternal(b *BTree, k, i int) {
	switch {
	case len(n.Children[i].Keys) > b.GetMinKeys():
		n.Keys[i] = n.deleteFromLeftChild(b, i)
	case len(n.Children[i+1].Keys) > b.GetMinKeys():
		n.Keys[i] = n.deleteFromRightChild(b, i+1)
	default:
		n.deleteAndMergeChildren(b, k, i)
	}
}

func (n *Node) delete(b *BTree, k int) {
	index := sort.Search(len(n.Keys), func(i int) bool {
		return n.Keys[i] >= k
	})

	if len(n.Children) == 0 && index < len(n.Keys) && n.Keys[index] == k {
		n.Keys = append(n.Keys[:index], n.Keys[index+1:]...)
		return
	}

	if index < len(n.Keys) && n.Keys[index] == k {
		n.deleteInternal(b, k, index)
	} else if len(n.Children[index].Keys) > b.GetMinKeys() {
		n.Children[index].delete(b, k)
	} else {
		child := n.Children[index]
		if index < len(n.Children)-1 && len(n.Children[index+1].Keys) > b.GetMinKeys() {
			child.Keys = append(child.Keys, n.Keys[index])
			n.Keys[index] = n.deleteFromRightChild(b, index+1)
			child.delete(b, k)
		} else if index > 0 && len(n.Children[index-1].Keys) > b.GetMinKeys() {
			child.Keys = append(child.Keys, n.Keys[index])
			n.Keys[index] = n.deleteFromLeftChild(b, index-1)
			child.delete(b, k)
		} else {
			if index == len(n.Children)-1 {
				n.mergeChildren(index - 1)
			} else {
				n.mergeChildren(index)
			}
			if index == 0 {
				n.Children[index].delete(b, k)
			} else {
				n.Children[index-1].delete(b, k)
			}
		}
	}
}

type BTree struct {
	BranchingFactor int
	Root            *Node
}

func (b *BTree) Insert(k int) {
	if b.Root == nil {
		b.Root = &Node{
			Keys:     []int{k},
			Children: make([]*Node, 0),
		}
		return
	}
	if len(b.Root.Keys) == b.GetMaxKeys() {
		newRoot := &Node{
			Keys:     make([]int, 0),
			Children: make([]*Node, 0),
		}
		newRoot.Children = append(newRoot.Children, b.Root)
		b.Root = newRoot
		newRoot.splitChild(0)
		newRoot.insert(b, k)
	} else {
		b.Root.insert(b, k)
	}
}

func (b *BTree) Delete(k int) {
	b.Root.delete(b, k)
	if len(b.Root.Keys) == 0 {
		b.Root = b.Root.Children[0]
	}
}

func (b *BTree) GetMaxKeys() int {
	return 2*b.BranchingFactor - 1
}

func (b *BTree) GetMinKeys() int {
	return b.BranchingFactor - 1
}

func (b *BTree) GetMaxChild() int {
	return 2 * b.BranchingFactor
}

// Used only for pretty print(debugging)
func (b *BTree) traverse() []prettyLine {
	lines := make([]prettyLine, 0)

	node := b.Root
	counters := [2]int{0, 1}
	depth := 0
	queue := make([]*Node, 0)
	queue = append(queue, node)
	for {
		if len(queue) == 0 {
			break
		}

		node = queue[0]
		if counters[0] == 0 {
			depth++
			counters[0] = counters[1]
			counters[1] = 0

			newLine := prettyLine{
				nodes: make([]*Node, 0),
				depth: depth,
			}
			lines = append(lines, newLine)
		}
		lines[len(lines)-1].nodes = append(lines[len(lines)-1].nodes, node)
		queue = queue[1:]
		if len(node.Children) > 0 {
			queue = append(queue, node.Children...)
			counters[1] += len(node.Children)
		}
		counters[0]--
	}
	return lines
}

// Used only for debugging
func (b *BTree) PrettyString() string {
	lines := b.traverse()

	lines = lines[1:]
	res := ""
	slices.Reverse(lines)
	maxLineLen := 0
	for _, line := range lines {

		s := line.generatePrettyLine()

		if maxLineLen != 0 {
			spaces := maxLineLen/2 - len(s)/2
			if spaces > 0 {
				s = strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces)
			}
		}
		if maxLineLen == 0 {
			maxLineLen = len(s)
		}
		res = s + "\n\n" + res
	}

	s := b.Root.String()
	spaces := maxLineLen/2 - len(s)/2
	s = strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces)
	res = s + "\n\n" + res

	return res
}

func NewBTree(b int) *BTree {
	return &BTree{
		BranchingFactor: b,
	}
}

// Used only for debugging(this is shitty implementation and should be removed in the future)
type prettyLine struct {
	nodes []*Node
	depth int
}

func (pl *prettyLine) generatePrettyLine() string {
	s := ""
	for _, n := range pl.nodes {
		childString := ""
		for _, c := range n.Children {
			childString += c.String()
		}

		if len(childString) > 0 {
			spaces := len(childString)/2 - len(n.String())/2
			s += strings.Repeat(" ", spaces) + n.String() + strings.Repeat(" ", spaces)
		} else {
			s += n.String()
		}
	}
	return s
}
