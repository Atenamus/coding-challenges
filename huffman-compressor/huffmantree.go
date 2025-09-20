package main

type HuffBaseNode interface {
	IsLeaf() bool
	Freq() int
}

type HuffLeafNode struct {
	char string
	freq int
}

func NewLeafNode(c string, f int) *HuffLeafNode {
	return &HuffLeafNode{
		char: c,
		freq: f,
	}
}

func (leaf *HuffLeafNode) Freq() int {
	return leaf.freq
}

func (leaf *HuffLeafNode) IsLeaf() bool {
	return true
}

func (leaf *HuffLeafNode) Char() string {
	return leaf.char
}

type HuffInternalNode struct {
	freq  int
	left  HuffBaseNode
	right HuffBaseNode
}

func NewInternalNode(f int, l HuffBaseNode, r HuffBaseNode) *HuffInternalNode {
	return &HuffInternalNode{
		freq:  f,
		left:  l,
		right: r,
	}
}

func (n *HuffInternalNode) Freq() int {
	return n.freq
}

func (n *HuffInternalNode) IsLeaf() bool {
	return false
}

func (n *HuffInternalNode) Left() HuffBaseNode {
	return n.left
}

func (n *HuffInternalNode) Right() HuffBaseNode {
	return n.right
}
