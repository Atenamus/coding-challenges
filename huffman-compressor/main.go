package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func PrintTree(node HuffBaseNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Choose the appropriate connector
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	if node.IsLeaf() {
		leaf := node.(*HuffLeafNode)
		// Display character with frequency, handle special characters
		char := leaf.Char()
		switch char {
		case " ":
			char = "SPACE"
		case "\n":
			char = "\\n"
		case "\t":
			char = "\\t"
		case "\r":
			char = "\\r"
		}
		fmt.Printf("%s%s'%s' (freq: %d)\n", prefix, connector, char, leaf.Freq())
	} else {
		internal := node.(*HuffInternalNode)
		fmt.Printf("%s%sInternal (freq: %d)\n", prefix, connector, internal.Freq())

		// Prepare prefix for children
		extension := "│   "
		if isLast {
			extension = "    "
		}
		newPrefix := prefix + extension

		// Print children - left first, then right
		PrintTree(internal.Left(), newPrefix, false)
		PrintTree(internal.Right(), newPrefix, true)
	}
}

func isEscapeChar(r rune) bool {
	switch r {
	case '\n', '\r', '\t', '\b', '\f', '\v', '\\':
		return true
	}
	return false
}

func main() {
	file, err := os.Open("test.txt")
	freq := make(map[string]int)
	nodes := []HuffBaseNode{}
	pq := &PriorityQueue{}

	heap.Init(pq)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading rune: %v", err)
		}
		if !isEscapeChar(r) {
			freq[string(r)]++
		}
	}

	for key, value := range freq {
		nodes = append(nodes, NewLeafNode(key, value))
	}

	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Freq() < nodes[j].Freq() })

	for _, val := range nodes {
		heap.Push(pq, val)
	}

	for pq.Len() > 1 {
		l := pq.Pop()
		r := pq.Pop()
		pq.Push(NewInternalNode(l.(HuffBaseNode).Freq()+r.(HuffBaseNode).Freq(), l.(HuffBaseNode), r.(HuffBaseNode)))
	}

	root := pq.Pop()
	PrintTree(root.(HuffBaseNode), "", true)
}
