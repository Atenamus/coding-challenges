package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func IsEscapeChar(r rune) bool {
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

func ReadFile(path string) (error, string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("File does not exist"), ""
	}
	text, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error reading file : %v", err), ""
	}
	return nil, string(text)
}

func DisplayStats(originalSize, compressedSize, headerSize uint32, codes map[string]string) {
	fmt.Printf("\n=== COMPRESSION STATISTICS ===\n")
	fmt.Printf("Original size: %d bytes\n", originalSize)
	fmt.Printf("Header size: %d bytes\n", headerSize)
	fmt.Printf("Compressed data: %d bytes\n", compressedSize)
	fmt.Printf("Total output: %d bytes\n", headerSize+compressedSize)
	ratio := float64(originalSize) / float64(headerSize+compressedSize)
	fmt.Printf("Compression ratio: %.2f:1\n", ratio)
	fmt.Printf("Space saved: %.1f%%\n", (1.0-float64(headerSize+compressedSize)/float64(originalSize))*100)
}

func main() {
	// inputFile := flag.String("input", "", "Input file to compress")
	freq := make(map[string]int)
	nodes := []HuffBaseNode{}
	pq := &PriorityQueue{}

	heap.Init(pq)
	flag.Parse()

	err, text := ReadFile("test.txt")
	if err != nil {
		log.Fatalf("Error while reading file : %v", err)
	}

	for _, char := range text {
		if !IsEscapeChar(char) {
			freq[string(char)]++
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
	codes := NewHuffmanCodesTable()

	if root, ok := root.(*HuffInternalNode); ok {
		codes = root.GenerateHuffmanCodes()
	}
	renderTable(codes)

	outFile, err := os.Create("test.huff")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	originalSize := uint32(len(text))
	headerSize, err := WriteFileHeader(writer, codes, originalSize)
	if err != nil {
		log.Fatalf("Error writing to header : %v", err)
	}
	compressedSize, err := EncodeText(text, codes, writer)
	if err != nil {
		log.Fatalf("Error encoding data: %v", err)
	}
	err = UpdateCompressedSize("test.huff", compressedSize)
	if err != nil {
		log.Fatalf("Error updating header: %v", err)
	}

	DisplayStats(originalSize, compressedSize, headerSize, codes)
}
