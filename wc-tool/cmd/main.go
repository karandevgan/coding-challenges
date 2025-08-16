package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
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

	reader := bufio.NewReader(f)
	inWord := false
	for {
		r, n, err := reader.ReadRune()
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
		tBytes += n
		tChars += 1
		if r == '\n' {
			tLines += 1
		}
		// check if char is whitespace
		if unicode.IsSpace(r) {
			if inWord {
				tWords += 1
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	if inWord {
		tWords += 1
	}

	var outputStrArr []string
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
