package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"unicode"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("No file provided")
		return
	}

	fileName := args[0]
	validJson := validateJSONFromFile(fileName)
	if validJson {
		fmt.Println("Valid JSON")
	} else {
		fmt.Println("Invalid JSON")
	}
}

func validateJSONFromFile(fileName string) bool {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return false
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	return validateJSON(reader)
}

func validateJSON(reader *bufio.Reader) bool {
	// Only allow a single top-level value that is either an object or an array.
	iterateOverContinuousWhitespace(reader)
	// Peek first non-whitespace rune to decide top-level type
	r, _, err := reader.ReadRune()
	if err != nil {
		return false
	}
	if r != '{' && r != '[' {
		// Top-level must be object or array
		return false
	}
	_ = reader.UnreadRune()
	var ok bool
	if r == '{' {
		ok = validateJSONObject(reader)
	} else {
		ok = validateJSONArray(reader)
	}
	if !ok {
		return false
	}
	iterateOverContinuousWhitespace(reader)
	_, _, err = reader.ReadRune()
	if err == io.EOF {
		return true
	}
	// Any non-whitespace after a complete value is invalid JSON
	return false
}

func validateJSONObject(reader *bufio.Reader) bool {
	/*
		Rules for a json object:
		1. JSON object must start with { and end with }
		2. JSON object must have a key value pair, and a comma after each key value pair
		3. Key must be enclosed in double quotes
	*/
	nextExpectedDelims := []rune{'{'}
	for {
		iterateOverContinuousWhitespace(reader)
		r, _, err := reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading file: %s\n", err)
			}
			return false
		}
		// After removing all the whitespace, first should be the end bracket or key
		if !slices.Contains(nextExpectedDelims, r) {
			return false
		}
		if r == '{' {
			nextExpectedDelims = []rune{'}', '"'}
		} else if r == '}' {
			return true
		} else if r == '"' {
			if !validateString(reader) {
				return false
			}
			nextExpectedDelims = []rune{':'}
		} else if r == ':' {
			if !validateJSONValue(reader) {
				return false
			}
			nextExpectedDelims = []rune{',', '}'}
		} else if r == ',' {
			nextExpectedDelims = []rune{'"'}
		} else {
			return false
		}
	}
}

func validateJSONArray(reader *bufio.Reader) bool {
	nextExpectedDelims := []rune{'['}
	for {
		iterateOverContinuousWhitespace(reader)
		r, _, err := reader.ReadRune()
		if err != nil {
			return false
		}
		if !slices.Contains(nextExpectedDelims, r) {
			return false
		}
		if r == '[' {
			iterateOverContinuousWhitespace(reader)
			r, _, err = reader.ReadRune()
			if err != nil {
				return false
			}
			if r == ']' {
				return true
			}
			_ = reader.UnreadRune()
			if !validateJSONValue(reader) {
				return false
			}
			nextExpectedDelims = []rune{']', ','}
		} else if r == ']' {
			return true
		} else if r == ',' {
			nextExpectedDelims = []rune{',', ']'}
			if !validateJSONValue(reader) {
				return false
			}
		}
	}
}

func validateJSONValue(reader *bufio.Reader) bool {
	/*
		Rules for a json value:
			1. JSON value can be a string, number, boolean, null, object, array
			2. JSON value can be enclosed in double quotes
	*/
	for {
		iterateOverContinuousWhitespace(reader)
		r, _, err := reader.ReadRune()
		if err != nil {
			return false
		}
		if r == '{' {
			_ = reader.UnreadRune()
			return validateJSONObject(reader)
		} else if r == '[' {
			_ = reader.UnreadRune()
			return validateJSONArray(reader)
		} else if r == '"' {
			return validateString(reader)
		} else if r == 'n' {
			return validateSequence(reader, []int32{'u', 'l', 'l'})
		} else if unicode.IsDigit(r) || r == '-' {
			_ = reader.UnreadRune()
			return validateNumber(reader)
		} else if r == 't' {
			return validateSequence(reader, []int32{'r', 'u', 'e'})
		} else if r == 'f' {
			return validateSequence(reader, []int32{'a', 'l', 's', 'e'})
		} else {
			return false
		}
	}
}

func validateNumber(reader *bufio.Reader) bool {
	isFirstDigit := true
	isLeadingZero := false
	signSeen := false
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			// EOF is valid only if we've consumed at least one digit
			if err == io.EOF {
				return !isFirstDigit
			}
			return false
		}
		if !unicode.IsDigit(r) {
			// Only a single leading '-' is allowed
			if r == '-' {
				if !isFirstDigit || signSeen {
					return false
				}
				signSeen = true
				continue
			}
			// No '+' allowed at any position in JSON numbers
			if r == '+' {
				return false
			}
			if isFirstDigit {
				return false
			}
			// If the number is a float, then it can have a single decimal point before the exponent
			if r == '.' {
				return validateDecimalPart(reader)
			} else if r == 'e' || r == 'E' {
				return validateExponentPart(reader)
			}
			_ = reader.UnreadRune()
			return true
		}
		if r == '0' && isFirstDigit {
			isLeadingZero = true
		}
		if isLeadingZero && !isFirstDigit {
			// No digits allowed after a leading zero unless followed by . or exponent (handled above)
			return false
		}
		isFirstDigit = false
	}
}

func validateDecimalPart(reader *bufio.Reader) bool {
	isFirstDigit := true
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// Valid only if we saw at least one digit after the decimal point
				return !isFirstDigit
			}
			return false
		}
		if !unicode.IsDigit(r) {
			// After decimal point, first character should be a digit
			if isFirstDigit {
				return false
			}
			if r == 'e' || r == 'E' {
				return validateExponentPart(reader)
			}
			_ = reader.UnreadRune()
			return true
		}
		isFirstDigit = false
	}
}

func validateExponentPart(reader *bufio.Reader) bool {
	isFirstDigit := true
	signSeen := false
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// Must have at least one digit in the exponent
				return !isFirstDigit
			}
			return false
		}
		if !unicode.IsDigit(r) {
			if r == '+' || r == '-' {
				// Allow at most one sign and only before the first digit
				if !isFirstDigit || signSeen {
					return false
				}
				signSeen = true
				continue
			}
			if isFirstDigit {
				return false
			}
			_ = reader.UnreadRune()
			return true
		}
		isFirstDigit = false
	}
}

func validateSequence(reader *bufio.Reader, seq []int32) bool {
	nextCharIndex := 0
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return false
		}
		if r != seq[nextCharIndex] {
			return false
		}
		nextCharIndex += 1
		if nextCharIndex == len(seq) {
			return true
		}
	}
}

func validateString(reader *bufio.Reader) bool {
	// Called after the opening double quote (") has been consumed.
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return false
		}
		// Closing quote (not escaped) ends the string
		if r == '"' {
			return true
		}
		// Unescaped control characters (U+0000 through U+001F) are not allowed in JSON strings
		if r < 0x20 {
			return false
		}
		if r == '\\' {
			// Validate escape sequence
			esc, _, err := reader.ReadRune()
			if err != nil {
				return false
			}
			switch esc {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				// Valid simple escapes
			case 'u':
				// Expect exactly four hex digits
				for i := 0; i < 4; i++ {
					h, _, err := reader.ReadRune()
					if err != nil {
						return false
					}
					if !isHexDigit(h) {
						return false
					}
				}
			default:
				return false
			}
			// Continue scanning characters after a valid escape sequence
			continue
		}
		// Any other rune is fine (already filtered out control chars)
	}
}

func iterateOverContinuousWhitespace(reader *bufio.Reader) {
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return
		}
		if unicode.IsSpace(r) {
			continue
		}
		_ = reader.UnreadRune()
		break
	}
}

// isHexDigit reports whether r is a hexadecimal digit.
func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}
