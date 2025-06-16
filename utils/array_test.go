package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBatchSliceInt(t *testing.T) {
	tests := []struct {
		name      string
		arr       []int
		batchSize int
		want      [][]int
	}{
		{
			name:      "empty slice",
			arr:       []int{},
			batchSize: 3,
			want:      [][]int{},
		},
		{
			name:      "batch size <= 0",
			arr:       []int{1, 2, 3},
			batchSize: 0,
			want:      [][]int{},
		},
		{
			name:      "batch size > len(arr)",
			arr:       []int{1, 2},
			batchSize: 5,
			want:      [][]int{{1, 2}},
		},
		{
			name:      "even batches",
			arr:       []int{1, 2, 3, 4},
			batchSize: 2,
			want:      [][]int{{1, 2}, {3, 4}},
		},
		{
			name:      "uneven batches",
			arr:       []int{1, 2, 3, 4, 5},
			batchSize: 2,
			want:      [][]int{{1, 2}, {3, 4}, {5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BatchSlice(tt.arr, tt.batchSize)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBatchSliceString(t *testing.T) {
	arr := []string{"a", "b", "c", "d", "e"}
	want := [][]string{{"a", "b"}, {"c", "d"}, {"e"}}
	got := BatchSlice(arr, 2)
	require.Equal(t, want, got)
}
