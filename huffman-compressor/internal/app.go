package internal

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Atenamus/compress-tool/internal/huffman"
)

type Application struct {
	logger *log.Logger
}

func NewApplication() *Application {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return &Application{
		logger: logger,
	}
}

func (a *Application) Run() {
	var (
		compress   = flag.Bool("compress", false, "compress file")
		decompress = flag.Bool("decompress", false, "decompress file")
		input      = flag.String("input", "", "input file path")
		output     = flag.String("output", "", "output file path")
	)

	flag.Parse()

	if err := validateArgs(*compress, *decompress, *input); err != nil {
		a.logger.Fatalf("Error : %v\n", err)
	}

	var err error
	if *compress {
		err = a.compress(*input, *output)
	} else {
		err = a.decompress(*input, *output)
	}

	if err != nil {
		a.logger.Fatalf("Failed : %v\n", err.Error())
	}

	a.logger.Println("Success")
}

func (a *Application) compress(input, output string) error {
	if output == "" {
		output = strings.TrimSuffix(input, filepath.Ext(input)) + ".huff"
	}
	stats, err := huffman.CompressFile(input, output)
	if err != nil {
		return err
	}

	displayStats(stats.OriginalSize, stats.CompressedSize, stats.HeaderSize)
	return nil
}

func (a *Application) decompress(input, output string) error {
	if output == "" {
		output = strings.TrimSuffix(input, ".huff") + "_extracted.txt"
	}

	if err := huffman.Decompress(input, output); err != nil {
		return err
	}
	return nil
}

func validateArgs(encode, decode bool, input string) error {
	if !encode && !decode {
		return fmt.Errorf("specify either -encode or -decode")
	}

	if encode && decode {
		return fmt.Errorf("cannot specify both -encode and -decode")
	}

	if input == "" {
		return fmt.Errorf("input cannot be empty")
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist")
	}

	return nil
}

func displayStats(originalSize, compressedSize, headerSize uint32) {
	fmt.Printf("\n=== COMPRESSION STATISTICS ===\n")
	fmt.Printf("Original size: %d bytes\n", originalSize)
	fmt.Printf("Header size: %d bytes\n", headerSize)
	fmt.Printf("Compressed data: %d bytes\n", compressedSize)
	fmt.Printf("Total output: %d bytes\n", headerSize+compressedSize)
	fmt.Printf("Space saved: %.1f%%\n", (1.0-float64(headerSize+compressedSize)/float64(originalSize))*100)
}
