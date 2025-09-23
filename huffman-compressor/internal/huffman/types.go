package huffman

type PriorityQueue []HuffBaseNode

type HuffmanCodes map[string]string

type FrquencyTable map[string]int

type HuffBaseNode interface {
	IsLeaf() bool
	Freq() int
}

type HuffLeafNode struct {
	char string
	freq int
}

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

type Statistics struct {
	OriginalSize   uint32
	CompressedSize uint32
	HeaderSize     uint32
	SpaceSaved     float64
}
