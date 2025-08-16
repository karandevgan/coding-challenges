package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	cBytes := flag.Bool("c", false, "count bytes")
	cLines := flag.Bool("l", false, "count lines")
	cWords := flag.Bool("w", false, "count words")
	mWords := flag.Bool("m", false, "count number of characters")
	flag.Parse()

	if !*cBytes && !*cLines && !*cWords && !*mWords {
		*cBytes = true
		*cLines = true
		*cWords = true
	}

	args := flag.Args()
	readStdin := false
	if len(args) < 1 {
		readStdin = true
	}
	const maxLine = 256 * 1024
	buf := make([]byte, maxLine)
	var f *os.File
	var fileName string
	var err error
	if readStdin {
		// read from stdin
		f = os.Stdin
		defer f.Close()
	} else {
		fileName = args[0]
		f, err = os.Open(fileName)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			return
		}
		defer f.Close()
	}

	tBytes := 0
	tLines := 0
	tWords := 0
	tChars := 0

	var outputStrArr []string
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			return
		}
		if n == 0 {
			break
		}
		buf = buf[:n]
		tBytes += n
		// count words
		wEnd := false
		var prevChar string
		for _, b := range buf {
			char := string(b)
			// check if char is whitespace
			if strings.Contains(" \t\n\r", char) && !wEnd {
				wEnd = true
				tWords += 1
			} else if !strings.Contains(" \t\n\r", char) {
				wEnd = false
			}
			if !strings.Contains("\n\r", char) {
				tChars += 1
			}
			if strings.Contains("\n\r", char) && prevChar != "\n" && char != "\r" {
				tChars += 1
				tLines += 1
			}
			prevChar = char
		}
	}

	if *cLines {
		outputStrArr = append(outputStrArr, strconv.Itoa(tLines))
	}
	if *cWords {
		outputStrArr = append(outputStrArr, strconv.Itoa(tWords))
	}
	if *cBytes {
		outputStrArr = append(outputStrArr, strconv.Itoa(tBytes))
	}
	if *mWords {
		outputStrArr = append(outputStrArr, strconv.Itoa(tChars))
	}
	outputStrArr = append(outputStrArr, fileName)
	outputStr := strings.Join(outputStrArr, " ")
	fmt.Println(outputStr)
}
