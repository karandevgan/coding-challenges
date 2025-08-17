package main

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestHuffmanTree(t *testing.T) {
	tests := []struct {
		name                string
		inputMap            map[rune]int32
		expectedTree        *HuffmanNode
		expectedLookupTable map[rune]lookupValue
	}{
		{
			name:                "test1",
			inputMap:            map[rune]int32{'C': 32, 'D': 42, 'E': 120, 'K': 7, 'L': 42, 'M': 24, 'U': 37, 'Z': 2},
			expectedTree:        getExpectedTreeTest1(),
			expectedLookupTable: getExpectedLookupTableTest1(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := buildHuffmanTree(tc.inputMap)
			if !compareNodes(actual, tc.expectedTree) {
				t.Fatalf("unexpected output: expectedTree %v, got %v", tc.expectedTree, actual)
			}
			actualLookupTable := make(map[rune]lookupValue)
			buildLookupTable(actual, actualLookupTable, 0, 0)
			if !maps.Equal(actualLookupTable, tc.expectedLookupTable) {
				t.Fatalf("unexpected output: expectedLookupTable %v, got %v", tc.expectedLookupTable, actualLookupTable)
			}
		})
	}
}

func TestCompression(t *testing.T) {
	tests := []struct {
		name                   string
		lookupTable            map[rune]lookupValue
		input                  string
		expectedCompressed     []uint32
		expectedLastBitsToRead uint8
	}{
		{
			name:                   "Test DEED",
			lookupTable:            getExpectedLookupTableTest1(),
			input:                  "DEED",
			expectedCompressed:     []uint32{2768240640},
			expectedLastBitsToRead: 24,
		},
		{
			name:                   "Test MUCK",
			lookupTable:            getExpectedLookupTableTest1(),
			input:                  "MUCK",
			expectedCompressed:     []uint32{4243537920},
			expectedLastBitsToRead: 14,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, aLastBitsToRead := compressString(tc.input, tc.lookupTable, uint32(0), uint8(32))
			if b := slices.Equal(tc.expectedCompressed, actual); !b {
				t.Fatalf("unexpected output: expectedCompressed %v, got %v", tc.expectedCompressed, actual)
			}
			if tc.expectedLastBitsToRead != aLastBitsToRead {
				t.Fatalf("unexpected output: expectedLastBitsToRead %v, got %v", tc.expectedLastBitsToRead, aLastBitsToRead)
			}
		})
	}
}

func TestDecompression(t *testing.T) {
	tests := []struct {
		name                 string
		huffmanTree          *HuffmanNode
		input                uint32
		lastBitsToRead       uint8
		expectedDecompressed string
	}{
		{
			name:                 "Test DEED",
			huffmanTree:          getExpectedTreeTest1(),
			input:                uint32(2768240640),
			lastBitsToRead:       8,
			expectedDecompressed: "DEED",
		},
		{
			name:                 "Test MUCK",
			huffmanTree:          getExpectedTreeTest1(),
			input:                uint32(4243537920),
			lastBitsToRead:       18,
			expectedDecompressed: "MUCK",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, runes, err := decompressString(tc.input, tc.lastBitsToRead, tc.huffmanTree, tc.huffmanTree)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.EqualFold(tc.expectedDecompressed, string(runes)) {
				t.Fatalf("unexpected output: expectedDecompressed %v, got %v", tc.expectedDecompressed, string(runes))
			}
		})
	}
}

func getExpectedLookupTableTest1() map[rune]lookupValue {
	lMap := make(map[rune]lookupValue)
	lMap['C'] = lookupValue{representation: 14, length: 4} // 1110
	lMap['D'] = lookupValue{representation: 5, length: 3}  // 101
	lMap['E'] = lookupValue{representation: 0, length: 1}  // 0
	lMap['K'] = lookupValue{representation: 61, length: 6} // 111101
	lMap['L'] = lookupValue{representation: 6, length: 3}  // 110
	lMap['M'] = lookupValue{representation: 31, length: 5} // 11111
	lMap['U'] = lookupValue{representation: 4, length: 3}  // 100
	lMap['Z'] = lookupValue{representation: 60, length: 6} // 111100
	return lMap
}

func getExpectedTreeTest1() *HuffmanNode {
	root := createHuffmanNode(0, 306, 0)
	root.left = createHuffmanNode('E', 120, 0)
	root.right = createHuffmanNode(0, 186, 1)
	root.right.left = createHuffmanNode(0, 79, 0)
	root.right.left.left = createHuffmanNode('U', 37, 0)
	root.right.left.right = createHuffmanNode('D', 42, 1)
	root.right.right = createHuffmanNode(0, 107, 1)
	root.right.right.left = createHuffmanNode('L', 42, 0)
	root.right.right.right = createHuffmanNode(0, 65, 1)
	root.right.right.right.left = createHuffmanNode('C', 32, 0)
	root.right.right.right.right = createHuffmanNode(0, 33, 1)
	root.right.right.right.right.left = createHuffmanNode(0, 9, 0)
	root.right.right.right.right.right = createHuffmanNode('M', 24, 1)
	root.right.right.right.right.left.left = createHuffmanNode('Z', 2, 0)
	root.right.right.right.right.left.right = createHuffmanNode('K', 7, 1)

	return root
}

func compareNodes(node1, node2 *HuffmanNode) bool {
	if node1 == nil && node2 == nil {
		return true
	}
	if node1 == nil || node2 == nil {
		return false
	}
	if node1.char != node2.char || node1.weight != node2.weight || node1.code != node2.code {
		return false
	}
	if compareNodes(node1.left, node2.left) {
		return compareNodes(node1.right, node2.right)
	}
	return false
}
