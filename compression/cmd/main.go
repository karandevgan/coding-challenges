package main

import (
	"flag"
	"fmt"
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
	lookupMap := make(map[rune]uint32)
	buildLookupTable(rootNode, lookupMap, 32, 0)
	fmt.Printf("Lookup Map: %v\n", lookupMap)
	//outFile := args[1]
	//fmt.Printf("Writing to file: %s\n", outFile)
	//os.WriteFile(outFile, []byte(fmt.Sprintf("%v", frequenciesFromFile)), 0644)
}
