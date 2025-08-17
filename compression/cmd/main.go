package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const magic = uint32(0x48554646)    // "HUFF" in ASCII
const remMagic = uint32(0x52454D42) // "REMB" in ASCII

func main() {
	compressFlag := flag.Bool("c", false, "compress")
	decompressFlag := flag.Bool("d", false, "decompress")
	flag.Parse()
	if !*compressFlag && !*decompressFlag {
		fmt.Println("No operation specified -c or -d")
		return
	}
	if *compressFlag && *decompressFlag {
		fmt.Println("Only one of -c or -d can be specified")
		return
	}
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Provide file to compress or decompress")
		return
	}
	inputFileName := args[0]
	if *compressFlag {
		compressFile(inputFileName)
	} else {
		decompressFile(inputFileName)
	}
}

func decompressFile(inputFileName string) {
	fmt.Printf("Processing File: %s\n", inputFileName)
	fileToRead, err := os.Open(inputFileName)
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		return
	}
	defer fileToRead.Close()

	// Get file size first
	fileInfo, err := fileToRead.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %s\n", err)
		return
	}
	fileSize := fileInfo.Size()

	// Read remaining bits from the very end (last byte)
	_, err = fileToRead.Seek(fileSize-1, 0)
	if err != nil {
		fmt.Printf("Error seeking to end of file: %s\n", err)
		return
	}

	var remainingBits uint8
	if err := binary.Read(fileToRead, binary.LittleEndian, &remainingBits); err != nil {
		fmt.Printf("Error reading remaining bits from file: %s\n", err)
		return
	}

	// Go back to beginning to read the header
	_, err = fileToRead.Seek(0, 0)
	if err != nil {
		fmt.Printf("Error seeking to beginning of file: %s\n", err)
		return
	}

	var wMagic uint32
	if err = binary.Read(fileToRead, binary.LittleEndian, &wMagic); err != nil {
		fmt.Printf("Error reading magic from file: %s\n", err)
	}
	if wMagic != magic {
		fmt.Printf("File does not have right format identifer: %d\n", wMagic)
		return
	}
	rootNode, err := readBinaryTree(fileToRead)
	if err != nil {
		fmt.Printf("Error reading tree from file: %s\n", err)
		return
	}

	extension := filepath.Ext(inputFileName)
	baseNameWithoutExt := strings.TrimSuffix(filepath.Base(inputFileName), extension)
	outputFileName := filepath.Join(filepath.Dir(inputFileName), baseNameWithoutExt+"_uncompressed.txt")
	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening output file: %s\n", err)
		return
	}
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	nxtNode := rootNode

	// Read the first word of compressed data
	var cur uint32
	if err = binary.Read(fileToRead, binary.LittleEndian, &cur); err != nil {
		if err != io.EOF {
			fmt.Printf("Error reading first compressed word: %s\n", err)
		}
		return
	}

	for {
		var next uint32
		if err = binary.Read(fileToRead, binary.LittleEndian, &next); err != nil {
			fmt.Printf("Error reading compressed data: %s\n", err)
			return
		}
		bitsToRead := uint8(32)
		if next == remMagic {
			if remainingBits == 32 {
				bitsToRead = 32
			} else {
				bitsToRead = 32 - remainingBits
			}
		}
		var decompressedRunes []rune
		nxtNode, decompressedRunes, err = decompressString(cur, bitsToRead, nxtNode, rootNode)
		if err != nil {
			fmt.Printf("Error decompressing string: %s\n", err)
			return
		}
		if len(decompressedRunes) > 0 {
			if _, err = writer.WriteString(string(decompressedRunes)); err != nil {
				fmt.Printf("Error writing decompressed data: %s\n", err)
				return
			}
		}
		if next == remMagic {
			break
		}
		cur = next
	}
}

func compressFile(inputFileName string) {
	fmt.Printf("Processing File: %s\n", inputFileName)
	frequenciesFromFile, err := getFrequenciesFromFile(inputFileName)
	if err != nil {
		fmt.Printf("Error creating frequency map from file: %s\n", err)
		return
	}
	//fmt.Printf("Frequency map: %v\n", frequenciesFromFile)
	rootNode := buildHuffmanTree(frequenciesFromFile)
	lookupMap := make(map[rune]lookupValue)
	buildLookupTable(rootNode, lookupMap, 0, 0)
	//fmt.Printf("Lookup Map: %v\n", lookupMap)

	fileToCompress, err := os.Open(inputFileName)
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		return
	}
	defer fileToCompress.Close()
	reader := bufio.NewReader(fileToCompress)
	b := make([]byte, 1024)

	extension := filepath.Ext(inputFileName)
	baseNameWithoutExt := strings.TrimSuffix(filepath.Base(inputFileName), extension)
	outputFileName := filepath.Join(filepath.Dir(inputFileName), baseNameWithoutExt+"_compressed.huff")
	fileToWrite, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error opening output file: %s\n", err)
		return
	}
	defer fileToWrite.Close()

	if err = binary.Write(fileToWrite, binary.LittleEndian, magic); err != nil {
		fmt.Printf("Error writing magic to file: %s\n", err)
		return
	}
	if err = writeBinaryTree(rootNode, fileToWrite); err != nil {
		fmt.Printf("Error writing tree to file: %s\n", err)
		return
	}

	seedUint32 := uint32(0)
	remainingBits := uint8(32)
	for {
		readCount, err := reader.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading from file: %s\n", err)
			return
		}
		if readCount == 0 {
			break
		}
		compressedData, tRemainingBits := compressString(string(b[:readCount]), lookupMap, seedUint32, remainingBits)
		if tRemainingBits > 0 {
			seedUint32 = compressedData[len(compressedData)-1] >> tRemainingBits
			compressedData = compressedData[:len(compressedData)-1]
			remainingBits = tRemainingBits
		} else {
			seedUint32 = uint32(0)
			remainingBits = 32
		}
		if err = writeCompressedData(compressedData, fileToWrite); err != nil {
			fmt.Printf("Error writing compressed data to file: %s\n", err)
			return
		}
	}

	if remainingBits < 32 {
		seedUint32 = seedUint32 << remainingBits
		if err = writeCompressedData([]uint32{seedUint32}, fileToWrite); err != nil {
			fmt.Printf("Error writing compressed data to file: %s\n", err)
			return
		}
	}
	if err = binary.Write(fileToWrite, binary.LittleEndian, remMagic); err != nil {
		fmt.Printf("Error writing remaining bit identifier file: %s\n", err)
	}
	if err = binary.Write(fileToWrite, binary.LittleEndian, remainingBits); err != nil {
		fmt.Printf("Error writing remaining bit to file: %s\n", err)
		return
	}
}

func writeCompressedData(compressedData []uint32, fileToWrite *os.File) error {
	for _, c := range compressedData {
		err := binary.Write(fileToWrite, binary.LittleEndian, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func readBinaryTree(fileToRead *os.File) (*HuffmanNode, error) {
	var marker byte
	var err error
	if err = binary.Read(fileToRead, binary.LittleEndian, &marker); err != nil {
		return nil, err
	}
	// if marker is 0, then we reached the end of the tree
	if marker == 0 {
		return nil, nil
	}
	node := &HuffmanNode{}
	if err = binary.Read(fileToRead, binary.LittleEndian, &node.char); err != nil {
		return nil, err
	}
	if err = binary.Read(fileToRead, binary.LittleEndian, &node.code); err != nil {
		return nil, err
	}
	if err = binary.Read(fileToRead, binary.LittleEndian, &node.weight); err != nil {
		return nil, err
	}
	if node.left, err = readBinaryTree(fileToRead); err != nil {
		return nil, err
	}
	if node.right, err = readBinaryTree(fileToRead); err != nil {
		return nil, err
	}
	return node, nil
}

func writeBinaryTree(node *HuffmanNode, fileToWrite *os.File) error {
	if node == nil {
		// Tell that next is nil
		return binary.Write(fileToWrite, binary.LittleEndian, byte(0))
	}
	// Tell that next is a node
	if err := binary.Write(fileToWrite, binary.LittleEndian, byte(1)); err != nil {
		return err
	}
	if err := binary.Write(fileToWrite, binary.LittleEndian, node.char); err != nil {
		return err
	}
	if err := binary.Write(fileToWrite, binary.LittleEndian, node.code); err != nil {
		return err
	}
	if err := binary.Write(fileToWrite, binary.LittleEndian, node.weight); err != nil {
		return err
	}
	if err := writeBinaryTree(node.left, fileToWrite); err != nil {
		return err
	}
	return writeBinaryTree(node.right, fileToWrite)
}
