package main

import (
	"maps"
	"testing"
)

func TestFrequencies(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected map[rune]int
		err      error
	}{
		{name: "test1", file: "../test_files/test.txt", expected: map[rune]int{'X': 333, 't': 223000}, err: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualMap, err := getFrequenciesFromFile(tc.file)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			keys := maps.Keys(tc.expected)
			for key := range keys {
				if tc.expected[key] != actualMap[key] {
					t.Fatalf("unexpected value for key %q: expected %d, got %d", key, tc.expected[key], actualMap[key])
				}
			}
		})
	}
}
