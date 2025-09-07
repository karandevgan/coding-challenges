package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"sort-tool/internal/sort"
)

func main() {
	uniq := flag.Bool("u", false, "get unique lines")
	// By default, use q sort
	useQSort := flag.Bool("qsort", true, "use merge sort")
	useMSort := flag.Bool("mergesort", false, "use merge sort")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return
	}
	fileName := args[len(args)-1]
	input, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
		return
	}
	defer input.Close()
	lines, err := readFromInput(input, *uniq)
	if err != nil {
		log.Fatalf("Error while reading from input: %v", err)
		return
	}
	var sorted []string
	if *useMSort {
		sorted = sort.MergeSort(lines)
	} else if *useQSort {
		sorted = sort.QuickSort(lines)
	}
	output := os.Stdout
	for _, line := range sorted {
		_, err := output.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Error while writing to output: %v", err)
			return
		}
	}
}

func readFromInput(input *os.File, uFilter bool) ([]string, error) {
	reader := bufio.NewReader(input)
	output := make([]string, 0)
	uniqMap := make(map[string]bool)
	for {
		sLine := ""
		line, p, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return output, err
		}
		sLine = string(line)
		for p {
			line, p, err = reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					return output, nil
				}
			}
			sLine += string(line)
		}
		if uFilter {
			if _, ok := uniqMap[sLine]; !ok {
				uniqMap[sLine] = true
			} else {
				continue
			}
		}
		output = append(output, sLine)
	}
	return output, nil
}
