package main

import (
	"bufio"
	"io"
	"os"
)

func getFrequenciesFromFile(file string) (map[rune]int32, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	return getFrequencies(reader)
}

func getFrequencies(reader *bufio.Reader) (freqMap map[rune]int32, err error) {
	freqMap = make(map[rune]int32)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		freqMap[r]++
	}
	return freqMap, nil
}
