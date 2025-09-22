package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type FileHeader struct {
	Magic          uint32
	HeaderSize     uint32
	CompressedSize uint32
	OriginalSize   uint32
	NumOfUniqChar  uint32
}
type PrefixTableEntry struct {
	Char       byte
	CodeLength uint8
	Code       uint32
}

func bitStringToUint32(bitString string) uint32 {
	var result uint32
	for i, bit := range bitString {
		if bit == '1' {
			result |= (1 << (31 - i))
		}
	}
	return result
}

func WriteFileHeader(writer io.Writer, codes HuffmanCodes, originalSize uint32) (uint32, error) {
	var buffer bytes.Buffer
	var tableEntries []PrefixTableEntry

	for char, code := range codes {
		entry := PrefixTableEntry{
			Char:       char[0],
			CodeLength: uint8(len(char)),
			Code:       bitStringToUint32(code),
		}
		tableEntries = append(tableEntries, entry)
	}
	sort.SliceStable(tableEntries, func(i, j int) bool {
		return tableEntries[i].Char < tableEntries[j].Char
	})

	for _, entry := range tableEntries {
		if err := binary.Write(&buffer, binary.LittleEndian, entry.Char); err != nil {
			return 0, fmt.Errorf("Error writing char to file :  %v\n", err)
		}
		if err := binary.Write(&buffer, binary.LittleEndian, entry.CodeLength); err != nil {
			return 0, fmt.Errorf("Error writing codeLength to file :  %v\n", err)
		}
		if err := binary.Write(&buffer, binary.LittleEndian, entry.Code); err != nil {
			return 0, fmt.Errorf("Error writing code to file :  %v\n", err)
		}
	}

	headerSize := uint32(binary.Size(FileHeader{})) + uint32(buffer.Len())

	header := FileHeader{
		Magic:          0x48554646,
		HeaderSize:     headerSize,
		CompressedSize: 0,
		OriginalSize:   originalSize,
		NumOfUniqChar:  uint32(len(tableEntries)),
	}

	if err := binary.Write(&buffer, binary.LittleEndian, header); err != nil {
		return 0, fmt.Errorf("Error writing header to file :  %v\n", err)
	}

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		return 0, fmt.Errorf("Error writing table to file :  %v\n", err)
	}

	return headerSize, nil
}

func EncodeText(text string, codes HuffmanCodes, writer io.Writer) (uint32, error) {
	var bitString strings.Builder
	for _, char := range text {
		if !IsEscapeChar(char) {
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
		return 0, fmt.Errorf("Error writing to file : %v\n", err)
	}

	byteWritten := uint32(1)
	for i := 0; i < bitString.Len(); i += 8 {
		var byteVal uint8
		for j := 0; j < 7; j++ {
			byteVal |= (1 << (7 - j))
		}
		if err := binary.Write(writer, binary.LittleEndian, byteVal); err != nil {
			return 0, fmt.Errorf("Error writing to file : %v\n", err)
		}
		byteWritten++
	}
	return byteWritten, nil
}

func UpdateCompressedSize(filename string, compressedSize uint32) error {
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
