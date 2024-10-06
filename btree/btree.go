package btree

import "sort"

type Node struct {
	Keys       []int
	Children   []*Node
	parentNode *Node
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

func (n *Node) maybeSplit(b *BTree) (*Node, bool) {

	if len(n.Keys) < b.GetMaxChild() {
		return n, false
	}

	leftNode := &Node{
		Keys:       make([]int, len(n.Keys)/2),
		Children:   make([]*Node, 0, len(n.Keys)/2+1),
		parentNode: n.parentNode,
	}
	copy(leftNode.Keys, n.Keys[:len(n.Keys)/2])

	rightNode := &Node{
		Keys:       make([]int, len(n.Keys)/2),
		Children:   make([]*Node, 0, len(n.Keys)/2+1),
		parentNode: n.parentNode,
	}
	copy(rightNode.Keys, n.Keys[len(n.Keys)/2+1:])

	if len(n.Children) > 0 {
		for i := 0; i < len(n.Children); i++ {
			if i < len(n.Children)/2 {
				leftNode.Children = append(leftNode.Children, n.Children[i])
				n.Children[i].parentNode = leftNode
			} else {
				rightNode.Children = append(rightNode.Children, n.Children[i])
				n.Children[i].parentNode = rightNode
			}
		}
	}

	leftoverKey := n.Keys[len(n.Keys)/2]

	if n.parentNode != nil {

		index := sort.Search(len(n.parentNode.Keys), func(i int) bool {
			return n.parentNode.Keys[i] > leftoverKey
		})
		n.parentNode.Keys = append(n.parentNode.Keys, 0)
		copy(n.parentNode.Keys[index+1:], n.parentNode.Keys[index:])
		n.parentNode.Keys[index] = leftoverKey

		oldChildren := n.parentNode.Children
		n.parentNode.Children = make([]*Node, len(n.parentNode.Children)+1)
		i := 0
		for ; i < index; i++ {
			n.parentNode.Children[i] = oldChildren[i]
		}
		n.parentNode.Children[i] = leftNode
		i++
		n.parentNode.Children[i] = rightNode
		i++
		for ; i < len(n.parentNode.Children); i++ {
			n.parentNode.Children[i] = oldChildren[i-1]
		}

		return n.parentNode.maybeSplit(b)
	} else {
		newNode := &Node{
			Keys:     []int{leftoverKey},
			Children: []*Node{leftNode, rightNode},
		}
		leftNode.parentNode = newNode
		rightNode.parentNode = newNode
		return newNode, true
	}
}

func (n *Node) Insert(b *BTree, k int) *Node {
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
		childNode.Insert(b, k)
	}

	//newNode, splitted := n.maybeSplit(b)
	//if splitted {
	//	return newNode
	//}

	return n
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
		newRoot.Insert(b, k)
	} else {
		b.Root.Insert(b, k)
	}
}

func (b *BTree) GetMaxKeys() int {
	return 2*b.BranchingFactor - 1
}

func (b *BTree) GetMaxChild() int {
	return 2 * b.BranchingFactor
}

func NewBTree(b int) *BTree {
	return &BTree{
		BranchingFactor: b,
	}
}
