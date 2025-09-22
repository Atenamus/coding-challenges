package main

type HuffmanCodes map[string]string

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

func NewHuffmanCodesTable() HuffmanCodes {
	return make(HuffmanCodes)
}

func (codes HuffmanCodes) GetCode(char string) (string, bool) {
	code, exists := codes[char]
	return code, exists
}

func (n *HuffInternalNode) GenerateHuffmanCodes() HuffmanCodes {
	codes := NewHuffmanCodesTable()
	traverse(n, "", codes)
	return codes
}

func traverse(node HuffBaseNode, currentCode string, codes HuffmanCodes) {
	if node == nil {
		return
	}

	if node.IsLeaf() {
		if leaf, ok := node.(*HuffLeafNode); ok {
			codes[leaf.Char()] = currentCode
		}
		return
	}

	if internal, ok := node.(*HuffInternalNode); ok {
		traverse(internal.Left(), currentCode+"0", codes)
		traverse(internal.Right(), currentCode+"1", codes)
	}
}
