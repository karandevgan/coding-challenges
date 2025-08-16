package main

import (
	"maps"
	"slices"
)

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

func buildLookupTable(node *HuffmanNode, lookupMap map[rune]uint32, bitToSet int, parentCode uint32) {
	// Build code by shifting parent's code left by 1 and setting LSB if current node's code is 1.
	// Use bitToSet as a depth marker to avoid introducing an extra leading zero at the root call.
	var x uint32
	if bitToSet == 32 { // root call
		x = parentCode
	} else {
		x = parentCode << 1
		if node.code == 1 {
			x |= 1
		}
	}
	if node.char != 0 {
		lookupMap[node.char] = x
	}
	if node.left != nil {
		buildLookupTable(node.left, lookupMap, bitToSet-1, x)
	}
	if node.right != nil {
		buildLookupTable(node.right, lookupMap, bitToSet-1, x)
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
