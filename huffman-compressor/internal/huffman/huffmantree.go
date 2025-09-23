package huffman

import (
	"container/heap"
	"sort"
)

func newLeafNode(c string, f int) *HuffLeafNode {
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

func newInternalNode(f int, l HuffBaseNode, r HuffBaseNode) *HuffInternalNode {
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

func newHuffmanCodesTable() HuffmanCodes {
	return make(HuffmanCodes)
}

func generateFrequencyTable(text string) FrquencyTable {
	freq := make(FrquencyTable)

	for _, char := range text {
		if !isEscapeChar(char) {
			freq[string(char)]++
		}
	}

	return freq
}

func buildHUffmanTree(freqTable FrquencyTable) *HuffInternalNode {
	pq := &PriorityQueue{}
	nodes := []HuffBaseNode{}

	heap.Init(pq)

	for key, value := range freqTable {
		nodes = append(nodes, newLeafNode(key, value))
	}

	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Freq() < nodes[j].Freq() })

	for _, val := range nodes {
		heap.Push(pq, val)
	}

	for pq.Len() > 1 {
		l := pq.Pop()
		r := pq.Pop()
		pq.Push(newInternalNode(l.(HuffBaseNode).Freq()+r.(HuffBaseNode).Freq(), l.(HuffBaseNode), r.(HuffBaseNode)))
	}

	root := pq.Pop()

	return root.(*HuffInternalNode)
}

func generateHuffmanCodes(n HuffBaseNode) HuffmanCodes {
	codes := newHuffmanCodesTable()
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
