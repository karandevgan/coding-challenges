package main

import (
	"maps"
	"testing"
)

func TestHuffmanTree(t *testing.T) {
	tests := []struct {
		name                string
		inputMap            map[rune]int
		expected            *HuffmanNode
		expectedLookupTable map[rune]uint32
	}{
		{name: "test1", inputMap: map[rune]int{'C': 32, 'D': 42, 'E': 120, 'K': 7, 'L': 42, 'M': 24, 'U': 37, 'Z': 2}, expected: getExpectedTreeTest1(), expectedLookupTable: getExpectedLookupTableTest1()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := buildHuffmanTree(tc.inputMap)
			if !compareNodes(actual, tc.expected) {
				t.Fatalf("unexpected output: expected %v, got %v", tc.expected, actual)
			}
			actualLookupTable := make(map[rune]uint32)
			buildLookupTable(actual, actualLookupTable, 32, 0)
			if !maps.Equal(actualLookupTable, tc.expectedLookupTable) {
				t.Fatalf("unexpected output: expected %v, got %v", tc.expectedLookupTable, actualLookupTable)
			}
		})
	}
}

func getExpectedLookupTableTest1() map[rune]uint32 {
	lMap := make(map[rune]uint32)
	lMap['C'] = 14 // 1110
	lMap['D'] = 5  // 101
	lMap['E'] = 0  // 0
	lMap['K'] = 61 // 111101
	lMap['L'] = 6  // 110
	lMap['M'] = 31 // 11111
	lMap['U'] = 4  // 100
	lMap['Z'] = 60 // 111100
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
