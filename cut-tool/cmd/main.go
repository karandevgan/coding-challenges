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
	col := 0
	delim := "\t"
	var err error
	args := os.Args[1:]
	for _, arg := range args {
		if strings.HasPrefix(arg, "-f") {
			if col, err = parseColumn(arg); err != nil {
				fmt.Println(err)
				return
			}
		} else if strings.HasPrefix(arg, "-d") {
			if delim, err = parseDelim(arg); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	if col == 0 {
		fmt.Println("No column specified")
		return
	}
	filename := args[len(args)-1]
	_ = readFileAndPrint(filename, col, delim)
}

func readFileAndPrint(filename string, col int, delim string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading file: %v\n", err)
				return err
			}
			break
		}
		s := internal.GetStringAtColumn(line, col, delim)
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

func parseColumn(arg string) (int, error) {
	if len(arg) < 3 {
		return 0, fmt.Errorf("-f should be followed by integer > 1\n")
	}
	suffix := arg[2:]
	if n, err := strconv.Atoi(suffix); err == nil {
		if n < 1 {
			return 0, fmt.Errorf("-f should be followed by integer > 1\n")
		}
		return n, nil
	}
	return 0, fmt.Errorf("Invalid column number: %s\n", suffix)
}
