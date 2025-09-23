package huffman

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

func uint32ToBitString(code uint32, length uint8) string {
	var result strings.Builder
	for i := int(length - 1); i >= 0; i-- {
		if (code>>i)&1 == 1 {
			result.WriteByte('1')
		} else {
			result.WriteByte('0')
		}
	}
	return result.String()
}

func Decompress(input, output string) error {
	file, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("Error opening file : %v\n", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	codes, header, err := readHeader(reader)
	if err != nil {
		return err
	}
	text, err := decodeText(reader, codes, header)
	if err != nil {
		return err
	}

	outFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("Error creating file : %v\n", err)
	}

	if _, err = outFile.WriteString(text); err != nil {
		return fmt.Errorf("Error writing to output file : %v\n", err)
	}

	defer outFile.Close()
	return nil
}

func readHeader(reader io.Reader) (HuffmanCodes, FileHeader, error) {
	var header FileHeader
	err := binary.Read(reader, binary.LittleEndian, &header)
	if err != nil {
		return nil, header, err
	}

	if header.Magic != 0x48554646 {
		return nil, header, fmt.Errorf("invalid file type")
	}

	codes := newHuffmanCodesTable()

	for i := uint32(0); i < header.NumOfUniqChar; i++ {
		var entry PrefixTableEntry

		err = binary.Read(reader, binary.LittleEndian, &entry.Char)
		if err != nil {
			return nil, header, err
		}

		err = binary.Read(reader, binary.LittleEndian, &entry.CodeLength)
		if err != nil {
			return nil, header, err
		}

		err = binary.Read(reader, binary.LittleEndian, &entry.Code)
		if err != nil {
			return nil, header, err
		}

		bitString := uint32ToBitString(entry.Code, entry.CodeLength)
		codes[string(entry.Char)] = bitString
	}

	return codes, header, nil
}

func decodeText(reader io.Reader, codes HuffmanCodes, header FileHeader) (string, error) {
	reversedCodes := make(map[string]string)
	for char, code := range codes {
		reversedCodes[code] = char
	}

	var padding uint8
	if err := binary.Read(reader, binary.LittleEndian, &padding); err != nil {
		return "", fmt.Errorf("error reading padding: %v", err)
	}

	compressedData := make([]byte, header.CompressedSize-1)
	_, err := io.ReadFull(reader, compressedData)
	if err != nil {
		return "", fmt.Errorf("error reading compressed data: %v", err)
	}

	var bitString strings.Builder
	for _, b := range compressedData {
		bitString.WriteString(fmt.Sprintf("%08b", b))
	}

	encodedBits := bitString.String()
	if len(encodedBits) > int(padding) {
		encodedBits = encodedBits[:len(encodedBits)-int(padding)]
	}

	var decodedText strings.Builder
	var currentCode strings.Builder
	for _, bit := range encodedBits {
		currentCode.WriteRune(bit)
		if char, found := reversedCodes[currentCode.String()]; found {
			decodedText.WriteString(char)
			currentCode.Reset()
			if uint32(decodedText.Len()) == header.OriginalSize {
				break
			}
		}
	}

	return decodedText.String(), nil
}
