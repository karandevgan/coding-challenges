package main

import (
	"fmt"
	"maps"
	"slices"
)

type lookupValue struct {
	representation uint
	length         uint
}
type HuffmanNode struct {
	char   rune
	weight int
	code   uint8
	left   *HuffmanNode
	right  *HuffmanNode
}

func createHuffmanNode(char rune, weight int, code uint8) *HuffmanNode {
	return &HuffmanNode{
		char:   char,
		weight: weight,
		code:   code,
	}
}

func buildLookupTable(node *HuffmanNode, lookupMap map[rune]lookupValue, depth int, parentCode uint) {
	// Build code by shifting parent's code left by 1 and setting LSB if current node's code is 1.
	// Use depth as a depth marker to avoid introducing an extra leading zero at the root call.
	var x uint
	if depth == 0 { // root call
		x = parentCode
	} else {
		x = parentCode << 1
		if node.code == 1 {
			x |= 1
		}
	}
	if node.char != 0 {
		lookupMap[node.char] = lookupValue{
			representation: x,
			length:         uint(depth),
		}
	}
	if node.left != nil {
		buildLookupTable(node.left, lookupMap, depth+1, x)
	}
	if node.right != nil {
		buildLookupTable(node.right, lookupMap, depth+1, x)
	}
}

func buildHuffmanTree(freqMap map[rune]int) *HuffmanNode {
	var nodes []*HuffmanNode
	for key := range maps.Keys(freqMap) {
		huffmanNode := &HuffmanNode{
			char:   key,
			weight: freqMap[key],
		}
		nodes = append(nodes, huffmanNode)
	}
	rootNode := createHuffmanTree(nodes)
	assignHuffmanCode(rootNode)
	return rootNode
}

func assignHuffmanCode(root *HuffmanNode) {
	if root == nil {
		return
	}
	if root.left != nil {
		root.left.code = 0
	}
	if root.right != nil {
		root.right.code = 1
	}
	assignHuffmanCode(root.left)
	assignHuffmanCode(root.right)
}

func createHuffmanTree(nodes []*HuffmanNode) *HuffmanNode {
	slices.SortFunc(nodes, func(e *HuffmanNode, e2 *HuffmanNode) int {
		d := e.weight - e2.weight
		if d == 0 {
			return int(e.char - e2.char)
		}
		return d
	})
	if len(nodes) == 0 {
		return nil
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	var newNodes []*HuffmanNode
	n1 := nodes[0]
	n2 := nodes[1]
	combinedNode := &HuffmanNode{
		weight: n1.weight + n2.weight,
		left:   n1,
		right:  n2,
	}
	newNodes = append(newNodes, combinedNode)
	newNodes = append(newNodes, nodes[2:]...)
	return createHuffmanTree(newNodes)
}

func compressString(input string, lookupMap map[rune]lookupValue) []uint32 {
	var compressed []uint32
	compressedUint32 := uint32(0)
	remainingBits := uint(32)
	for i, c := range input {
		lValue := lookupMap[c]
		v := uint32(lValue.representation)
		n := lValue.length
		if n > remainingBits {
			diff := n - remainingBits
			compressedUint32 = compressedUint32<<remainingBits | v>>diff
			compressed = append(compressed, compressedUint32)
			compressedUint32 = v & ((1 << diff) - 1)
			remainingBits = 32 - diff
		} else {
			// shift compressedUint32 left by n bits and OR u to it.
			compressedUint32 = compressedUint32<<n | v
			remainingBits -= n
		}
		// Need only to be added if the compressed string is not a multiple of 32 bits and have remaining bits.
		if i == len(input)-1 && remainingBits < 32 {
			compressedUint32 = compressedUint32 << remainingBits
			compressed = append(compressed, compressedUint32)
		}
	}
	return compressed
}

func decompressString(input []uint32, root *HuffmanNode) (nxtNode *HuffmanNode, output []rune, err error) {
	lIndex := 0
	nxtNode = root
	for _, w := range input {
		nxtNode, err = getNextNode(w, nxtNode)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting next node: %v", err)
		}
		if nxtNode.left == nil && nxtNode.right == nil {
			output = append(output, nxtNode.char)
			nxtNode = root
		}
		lIndex++
	}
	return nxtNode, output, nil
}

func getNextNode(w uint32, node *HuffmanNode) (nxtNode *HuffmanNode, err error) {
	if w == 0 {
		return node.left, nil
	} else {
		return node.right, nil
	}
}
