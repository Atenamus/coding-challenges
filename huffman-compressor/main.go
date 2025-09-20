package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func isEscapeChar(r rune) bool {
	switch r {
	case '\n', '\r', '\t', '\b', '\f', '\v', '\\':
		return true
	}
	return false
}

func main() {
	file, err := os.Open("example.txt")
	freq := make(map[string]int)
	nodes := make([]HuffLeafNode, 1)
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
		nodes = append(nodes, *NewLeafNode(key, value))
	}
	fmt.Println(nodes)
}
