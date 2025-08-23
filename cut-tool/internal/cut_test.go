package internal

import "testing"

func TestGetStringAtColumn(t *testing.T) {
	tests := []struct {
		name string
		s    string
		col  int
		sep  string
		want string
	}{
		{
			name: "At Index 2 with tab",
			s:    "this	is	a	test",
			col:  2,
			sep:  "\t",
			want: "is",
		},
		{
			name: "At Index 2 with comma",
			s:    "this,is a,test",
			col:  2,
			sep:  ",",
			want: "is a",
		},
		{
			name: "At Index 1",
			s:    "this	is	a	test",
			col:  1,
			sep:  "\t",
			want: "this",
		},
		{
			name: "At Index 5",
			s:    "this	is	a	test",
			col:  5,
			sep:  "\t",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetStringAtColumn(tc.s, tc.col, tc.sep)
			if actual != tc.want {
				t.Fatalf("unexpected output: expected %v, got %v", tc.want, actual)
			}
		})
	}
}
