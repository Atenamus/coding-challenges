package main

import (
	"container/heap"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func isEscapeChar(r rune) bool {
	switch r {
	case '\n', '\r', '\t', '\b', '\f', '\v', '\\':
		return true
	}
	return false
}

func renderTable(codes HuffmanCodes) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Char", "Freq")
	keys := []string{}
	for char := range codes {
		keys = append(keys, char)
	}
	for _, val := range keys {
		char := val
		if char == " " {
			char = "SPACE"
		}
		table.Append([]string{char, codes[val]})
	}
	table.Render()
}

func readFile(path string) (error, string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("File does not exist"), ""
	}
	text, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Eroor reading file : %v", err), ""
	}
	return nil, string(text)
}

func main() {
	freq := make(map[string]int)
	nodes := []HuffBaseNode{}
	pq := &PriorityQueue{}

	heap.Init(pq)

	err, text := readFile("example.txt")
	if err != nil {
		fmt.Println(err)
	}

	for _, char := range text {
		if !isEscapeChar(char) {
			freq[string(char)]++
		}
	}
	
	fmt.Println(freq)

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
	codes := NewHuffmanCodesTable()

	if root, ok := root.(*HuffInternalNode); ok {
		codes = root.GenerateHuffmanCodes()
	}
	renderTable(codes)
}
