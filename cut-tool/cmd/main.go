package main

import (
	"bufio"
	"cut-tool/internal"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	pCol := 0
	var oCols []int
	delim := "\t"
	var err error
	args := os.Args[1:]
	isLastIndexAFlag := false
	for i, arg := range args {
		if strings.HasPrefix(arg, "-f") {
			if pCol, oCols, err = parseColumns(arg); err != nil {
				fmt.Println(err)
				return
			}
			if i == len(args)-1 {
				isLastIndexAFlag = true
			}
		} else if strings.HasPrefix(arg, "-d") {
			if delim, err = parseDelim(arg); err != nil {
				fmt.Println(err)
				return
			}
			if i == len(args)-1 {
				isLastIndexAFlag = true
			}
		}
	}
	if pCol == 0 {
		fmt.Println("No column specified")
		return
	}
	lArg := args[len(args)-1]
	var reader *bufio.Reader
	if !isLastIndexAFlag && lArg != "-" {
		file, err := os.Open(lArg)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}
	_ = readFromReaderAndCut(reader, pCol, oCols, delim)
}

func readFromReaderAndCut(reader *bufio.Reader, col int, oCols []int, delim string) error {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading file: %v\n", err)
				return err
			}
			if len(line) > 0 {
				s := internal.GetStringAtColumn(line, delim, col, oCols...)
				fmt.Println(s)
			}
			break
		}
		s := internal.GetStringAtColumn(line, delim, col, oCols...)
		fmt.Println(s)
	}
	return nil
}

func parseDelim(arg string) (string, error) {
	if len(arg) < 3 {
		return "", fmt.Errorf("-d should be followed by delimiter\n")
	}
	return arg[2:], nil
}

func parseColumns(arg string) (int, []int, error) {
	if len(arg) < 3 {
		return 0, nil, fmt.Errorf("-f should be followed by integer > 1\n")
	}
	suffix := arg[2:]
	pCol := 0
	var oCols []int
	var cols []string
	if strings.Contains(suffix, ",") {
		cols = strings.Split(suffix, ",")
	} else if strings.Contains(suffix, " ") {
		cols = strings.Split(suffix, " ")
	} else {
		cols = []string{suffix}
	}
	if len(cols) == 0 {
		return 0, nil, fmt.Errorf("-f should be followed by integer > 1\n")
	}
	for i, col := range cols {
		if n, err := strconv.Atoi(strings.TrimSpace(col)); err == nil {
			if n < 1 {
				return 0, nil, fmt.Errorf("-f should be followed by integer > 1\n")
			}
			if i == 0 {
				pCol = n
			} else {
				oCols = append(oCols, n)
			}
		} else {
			return 0, nil, fmt.Errorf("Invalid column number: %s\n", col)
		}
	}

	return pCol, oCols, nil
}
