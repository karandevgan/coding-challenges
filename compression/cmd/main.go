package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Provide input and output file names")
		return
	}
	inputFileName := args[0]
	fmt.Printf("Compressing File: %s\n", inputFileName)
	frequenciesFromFile, err := getFrequenciesFromFile(inputFileName)
	if err != nil {
		fmt.Printf("Error creating frequency map from file: %s\n", err)
		return
	}
	//fmt.Printf("Frequency map: %v\n", frequenciesFromFile)
	rootNode := buildHuffmanTree(frequenciesFromFile)
	lookupMap := make(map[rune]lookupValue)
	buildLookupTable(rootNode, lookupMap, 0, 0)
	fmt.Printf("Lookup Map: %v\n", lookupMap)

	fileToCompress, err := os.Open(inputFileName)
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		return
	}
	defer fileToCompress.Close()
	reader := bufio.NewReader(fileToCompress)
	b := make([]byte, 1024)

	outputFileName := args[1]
	fileToWrite, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error opening output file: %s\n", err)
		return
	}
	defer fileToWrite.Close()
	writer := bufio.NewWriter(fileToWrite)

	seedUint32 := uint32(0)
	addSeedInData := false
	remainingBits := uint(32)
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
		compressedData, tRemainingBits := compressString(string(b), lookupMap, seedUint32, remainingBits)
		if tRemainingBits > 0 {
			seedUint32 = compressedData[len(compressedData)-1]
			seedUint32 = seedUint32 >> tRemainingBits
			compressedData = compressedData[:len(compressedData)-1]
			for _, c := range compressedData {
				err = writer.WriteByte(byte(c))
				if err != nil {
					fmt.Printf("Error writing to file: %s\n", err)
					return
				}
			}
			remainingBits = tRemainingBits
			addSeedInData = true
		} else {
			seedUint32 = uint32(0)
			addSeedInData = false
			remainingBits = 32
		}
	}

	if addSeedInData {
		// Add seed in final data with remaining bits
		fmt.Printf("Seed Data to be added")
		err = writer.WriteByte(byte(seedUint32))
		if err != nil {
			fmt.Printf("Error writing to file: %s\n", err)
			return
		}
		err = writer.WriteByte(byte(remainingBits))
		if err != nil {
			fmt.Printf("Error writing to file: %s\n", err)
			return
		}
	}

	//outFile := args[1]
	//fmt.Printf("Writing to file: %s\n", outFile)
	//os.WriteFile(outFile, []byte(fmt.Sprintf("%v", frequenciesFromFile)), 0644)
}
