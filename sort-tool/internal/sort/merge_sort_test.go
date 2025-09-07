package sort

import (
	"slices"
	"strings"
	"testing"
)

func TestMergeSort(t *testing.T) {
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
			got := MergeSort[string](tc.input)
			if !slices.Equal(tc.want, got) {
				t.Fatalf("unexpected output: expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestMergeSortInt(t *testing.T) {
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
			got := MergeSort[int](tc.input)
			if !slices.Equal(tc.want, got) {
				t.Fatalf("unexpected output: expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name      string
		l1        []string
		l2        []string
		want      []string
		expectErr bool
		err       error
	}{
		{
			name:      "l1 and l2 are empty",
			l1:        nil,
			l2:        make([]string, 0),
			want:      make([]string, 0),
			expectErr: false, err: nil,
		},
		{
			name:      "l1 is empty",
			l1:        make([]string, 0),
			l2:        []string{"Five", "Four", "One", "Three", "Two"},
			want:      []string{"Five", "Four", "One", "Three", "Two"},
			expectErr: false,
			err:       nil,
		},
		{
			name:      "l2 is empty",
			l1:        []string{"Five", "Four", "One", "Three", "Two"},
			l2:        make([]string, 0),
			want:      []string{"Five", "Four", "One", "Three", "Two"},
			expectErr: false,
			err:       nil,
		},
		{
			name:      "l1 and l2 are not empty",
			l1:        []string{"Five", "Four", "One", "Three", "Two"},
			l2:        []string{"Eight", "Nine", "Seven", "Six", "Ten"},
			want:      []string{"Eight", "Five", "Four", "Nine", "One", "Seven", "Six", "Ten", "Three", "Two"},
			expectErr: false,
			err:       nil,
		},
		{
			name:      "l1 size is > l2",
			l1:        []string{"Five", "Four", "One", "Seven", "Six", "Three", "Two"},
			l2:        []string{"Eight", "Nine", "Ten"},
			want:      []string{"Eight", "Five", "Four", "Nine", "One", "Seven", "Six", "Ten", "Three", "Two"},
			expectErr: false,
			err:       nil,
		},
		{
			name:      "l1 size is < l2",
			l1:        []string{"Eight", "Nine", "Ten"},
			l2:        []string{"Five", "Four", "One", "Seven", "Six", "Three", "Two"},
			want:      []string{"Eight", "Five", "Four", "Nine", "One", "Seven", "Six", "Ten", "Three", "Two"},
			expectErr: false,
			err:       nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := merge(tc.l1, tc.l2, strings.Compare)
			if !slices.Equal(tc.want, got) {
				t.Fatalf("unexpected output: expected %v, got %v", tc.want, got)
			}
		})
	}
}
