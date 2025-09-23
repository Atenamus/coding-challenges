package huffman

import (
	"fmt"
	"os"
)

func isEscapeChar(r rune) bool {
	switch r {
	case '\n', '\r', '\t', '\b', '\f', '\v', '\\':
		return true
	}
	return false
}

func readFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("File does not exist")
	}
	text, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Error reading file : %v", err)
	}
	return string(text), nil
}
