package sort

import (
	"slices"
	"testing"
)

func TestQuickSort(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		want      []string
		expectErr bool
		err       error
	}{
		{
			name:      "Input is empty",
			want:      nil,
			expectErr: false,
			err:       nil,
		},
		{
			name:      "Input is single value",
			input:     []string{"Five"},
			want:      []string{"Five"},
			expectErr: false,
			err:       nil,
		},
		{
			name:  "Input is multiple values",
			input: []string{"One", "Two", "Three", "Four", "Five"},
			want:  []string{"Five", "Four", "One", "Three", "Two"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := QuickSort(tc.input)
			if !slices.Equal(tc.want, got) {
				t.Fatalf("unexpected output: expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestQuickSortInt(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		want      []int
		expectErr bool
		err       error
	}{
		{
			name:  "Input is int",
			input: []int{-1, 4, 6, 2, 1, 3, 6, 5},
			want:  []int{-1, 1, 2, 3, 4, 5, 6, 6},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := QuickSort[int](tc.input)
			if !slices.Equal(tc.want, got) {
				t.Fatalf("unexpected output: expected %v, got %v", tc.want, got)
			}
		})
	}
}
