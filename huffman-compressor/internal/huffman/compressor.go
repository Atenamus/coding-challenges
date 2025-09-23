package huffman

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func bitStringToUint32(bitString string) uint32 {
	var result uint32
	for _, bit := range bitString {
		result <<= 1
		if bit == '1' {
			result |= 1
		}
	}
	return result
}

func CompressFile(input, output string) (*Statistics, error) {
	text, err := readFile(input)
	if err != nil {
		return nil, err
	}
	freqTable := generateFrequencyTable(text)
	tree := buildHUffmanTree(freqTable)
	codes := generateHuffmanCodes(tree)

	outFile, err := os.Create(output)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	originalSize := uint32(len(text))
	headerSize, err := writeFileHeader(writer, codes, originalSize)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nHeader size %d\n", headerSize)

	bytesWritten, err := encodeText(writer, text, codes)
	if err != nil {
		return nil, err
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}

	if err := updateCompressedSize(output, bytesWritten); err != nil {
		return nil, err
	}

	stats := &Statistics{
		OriginalSize:   originalSize,
		CompressedSize: bytesWritten,
		HeaderSize:     headerSize,
		SpaceSaved:     (1.0 - float64(headerSize+bytesWritten)/float64(originalSize)) * 100,
	}
	fmt.Println("Compression complete")
	return stats, nil
}

func writeFileHeader(writer io.Writer, codes HuffmanCodes, originalSize uint32) (uint32, error) {
	var tableBuffer bytes.Buffer
	var tableEntries []PrefixTableEntry

	for char, code := range codes {
		entry := PrefixTableEntry{
			Char:       char[0],
			CodeLength: uint8(len(code)),
			Code:       bitStringToUint32(code),
		}
		tableEntries = append(tableEntries, entry)
	}
	sort.SliceStable(tableEntries, func(i, j int) bool {
		return tableEntries[i].Char < tableEntries[j].Char
	})

	for _, entry := range tableEntries {
		if err := binary.Write(&tableBuffer, binary.LittleEndian, entry.Char); err != nil {
			return 0, fmt.Errorf("Error writing char to file :  %v\n", err)
		}
		if err := binary.Write(&tableBuffer, binary.LittleEndian, entry.CodeLength); err != nil {
			return 0, fmt.Errorf("Error writing codeLength to file :  %v\n", err)
		}
		if err := binary.Write(&tableBuffer, binary.LittleEndian, entry.Code); err != nil {
			return 0, fmt.Errorf("Error writing code to file :  %v\n", err)
		}
	}

	headerSize := uint32(binary.Size(FileHeader{})) + uint32(tableBuffer.Len())

	header := FileHeader{
		Magic:          0x48554646,
		HeaderSize:     headerSize,
		CompressedSize: 0,
		OriginalSize:   originalSize,
		NumOfUniqChar:  uint32(len(tableEntries)),
	}

	if err := binary.Write(writer, binary.LittleEndian, header); err != nil {
		return 0, fmt.Errorf("Error writing header to file :  %v\n", err)
	}

	if _, err := writer.Write(tableBuffer.Bytes()); err != nil {
		return 0, fmt.Errorf("Error writing table to file :  %v\n", err)
	}

	return headerSize, nil
}

func encodeText(writer io.Writer, text string, codes HuffmanCodes) (uint32, error) {
	var bitString strings.Builder
	for _, char := range text {
		if !isEscapeChar(char) {
			if code, exists := codes[string(char)]; exists {
				bitString.WriteString(code)
			} else {
				return 0, fmt.Errorf("Cannot find code for :%c", char)
			}
		}
	}

	bitStr := bitString.String()

	padding := (8 - (len(bitStr) % 8)) % 8
	bitStr += strings.Repeat("0", padding)

	if err := binary.Write(writer, binary.LittleEndian, uint8(padding)); err != nil {
		return 0, fmt.Errorf("Error writing padding: %v", err)
	}
	bytesWritten := uint32(1)

	for i := 0; i < len(bitStr); i += 8 {
		var byteVal uint8 = 0
		for j := range 8 {
			if bitStr[i+j] == '1' {
				byteVal |= 1 << (7 - j)
			}
		}
		if err := binary.Write(writer, binary.LittleEndian, byteVal); err != nil {
			return 0, fmt.Errorf("Error writing byte: %v", err)
		}
		bytesWritten++
	}

	return bytesWritten, nil
}

func updateCompressedSize(filename string, compressedSize uint32) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(8, 0)
	if err != nil {
		return err
	}

	return binary.Write(file, binary.LittleEndian, compressedSize)
}
